# Audit de JustDeploy — 1 novembre 2026

## Périmètre et méthode

Revue du code local : API Go/Echo, authentification, GitHub, Docker/Traefik, SQLite, interface React, installation et CI. Analyse statique et contrôles locaux, sans lancer le serveur ni modifier des conteneurs ou des données de production. Les impacts décrits sont déduits du code ; ils ne constituent pas des exploits exécutés sur un serveur.

Les niveaux sont des priorités de correction, sans score CVSS. L'exposition réseau réelle, les permissions du serveur et les sauvegardes externes n'ont pas été inspectées. Aucun correctif applicatif n'est inclus dans cet audit.

## Vérifications exécutées

| Contrôle | Résultat |
| --- | --- |
| `go test ./...` | Réussi. Un seul fichier de tests, dans `internal/utils` ; aucun test des déploiements, de l'API ou de l'authentification. |
| `go vet ./...` | Échec : deux balises JSON invalides dans `internal/domain/settings.go`. |
| `web`: `npm run lint` | Échec : 7 erreurs et 7 avertissements. |
| `web`: `tsc --noEmit` | Échec : 18 diagnostics, dont modules manquants et types incompatibles. |
| Build Vite, sortie dans `/tmp/justdeploy-audit-web` | Réussi. Le script de build n'exécute pas le contrôle TypeScript. |
| `npm audit --json` | Non effectué par npm : absence de lockfile npm (le projet utilise Bun). Aucun lockfile ajouté pour l'audit. |

Pas de scan CVE Go/frontend complet, ni de scan des images Docker. Aucune absence de vulnérabilités dans les dépendances ne peut être conclue. Build de l'image Docker et validation dans un navigateur non exécutés.

## Constats prioritaires

### 1. Critique — stockage des bases non rattaché lors de la recréation

**Sources :** `internal/adapter/docker.go:353`, `internal/application/deploy-application.go:111`.

`RunImageWithTLS` arrête et supprime le conteneur, puis le recrée avec un `HostConfig` vide : aucun volume nommé ni bind de données n'est défini. Les volumes anonymes éventuellement créés par l'image ne sont pas rattachés au nouveau conteneur. Une base peut donc redémarrer vide après un déploiement. Les anciennes données peuvent subsister dans un volume orphelin, sans être récupérées automatiquement ; les données stockées dans la couche du conteneur supprimé sont perdues.

**Correction :** volumes nommés stables par service, chemins de données explicites par moteur, sauvegarde/restauration vérifiée. Ne pas recréer les bases lors d'un déploiement applicatif sans changement les concernant.

### 2. Élevée — webhook GitHub sans authentification

**Sources :** `internal/api/routes.go:23`, `internal/api/handlers/post-github-event.go`, `internal/application/manage-github-event.go`, `internal/application/re-deploy-application.go`.

L'endpoint public accepte du JSON sans vérifier `X-Hub-Signature-256`. La seule condition est `ref == refs/heads/main`. Un appel externe peut déclencher un redéploiement si l'intégration GitHub est utilisable. La recherche utilise seulement le nom court du dépôt ; même sans service correspondant, le code finit par appeler `DeployApplication`, qui traite tous les services. Pas de contrôle du type d'événement ni de déduplication des livraisons.

**Impact :** reconstructions non autorisées, interruptions, consommation de ressources ; combinaison particulièrement grave avec le constat 1. Cela ne prouve pas qu'un attaquant peut substituer arbitrairement le code du dépôt.

**Correction :** vérifier le HMAC du corps brut avec le secret webhook, comparer en temps constant, vérifier événement, installation, dépôt complet et branche, dédupliquer les livraisons et mettre le travail en file.

### 3. Élevée — secrets GitHub et applicatifs dans les logs

**Sources :** `internal/adapter/git.go:28`, `internal/application/deploy-application.go:39`, `internal/application/create-app-from-repo.go:56`.

L'URL de clonage contenant le jeton GitHub est imprimée. Les structures de services et Compose sont aussi imprimées avec leurs variables d'environnement. L'URL authentifiée peut en outre rester dans `.git/config` du clone. Le contexte de build est archivé sans exclusions dans `makeTar`, ce qui transmet aussi `.git` et d'éventuels fichiers sensibles au builder ; leur incorporation dans l'image dépend du Dockerfile.

**Correction :** supprimer et masquer les secrets dans les logs, utiliser un mécanisme d'identification Git sans jeton dans l'URL persistée, retirer les credentials du remote et exclure `.git` et les secrets du contexte. Révoquer les jetons exposés et renouveler les secrets concernés si ces parcours ont déjà été utilisés.

### 4. Élevée — un service déclaré privé reste routé publiquement

**Sources :** `internal/application/deploy-application.go:87`, `internal/adapter/docker.go:357`, `web/src/components/ServiceDeploySettings.tsx:51`.

Le backend construit systématiquement un domaine pour les applications GitHub, même avec `IsExposed=false`. Dès qu'un domaine est fourni, `RunImageWithTLS` ajoute les labels Traefik. Le réglage de visibilité de l'interface n'est donc pas une barrière effective.

**Correction :** n'ajouter les labels et l'URL que lorsque l'exposition est autorisée ; retirer la route existante lorsque l'option est désactivée. Tester les deux transitions.

### 5. Élevée, conditionnelle — prise de contrôle de la première configuration

**Sources :** `internal/api/routes.go:25`, `internal/api/api.go`, `internal/application/save-initial-setup.go:11`, `internal/api/handlers/post-save-initial-setup.go`.

Le serveur écoute sur toutes les interfaces et `/api/setup` est public tant qu'aucun email administrateur n'est enregistré. Quiconque atteint une installation vierge peut devenir administrateur. Le handler n'impose aucune validation de l'email ou de la robustesse du mot de passe. Un email vide laisse le verrou de configuration ouvert. Le contrôle puis l'enregistrement ne sont pas atomiques : des requêtes simultanées peuvent toutes franchir le contrôle initial.

**Correction :** secret d'initialisation à usage unique, configuration locale ou accès réseau restreint pendant l'amorçage ; validation serveur et transaction conditionnelle pour créer l'administrateur une seule fois.

### 6. Élevée, conditionnelle — authentification servie en HTTP

**Sources :** `internal/api/api.go`, `web/src/services/createGithubApp.ts:4`, `internal/application/jwt.go`, `web/src/services/authStorage.ts`.

Le serveur d'administration démarre en HTTP sur 5915, et les URLs du manifeste GitHub imposent `http://IP:5915`. Sans terminaison TLS externe, mots de passe, JWT et codes GitHub transitent sans chiffrement. L'option HTTPS des services ne protège pas automatiquement cette API. Le JWT est conservé dans localStorage pendant une validité de 30 jours ; aucune révocation de session n'a été trouvée dans les routes.

**Correction :** origine publique HTTPS configurable pour l'administration et GitHub, durée de session réduite et révocation. Le stockage local augmente l'impact d'une éventuelle XSS ; aucune XSS exploitable n'est démontrée ici.

### 7. Élevée — isolation insuffisante entre applications

**Sources :** `internal/adapter/docker.go:364`, `internal/application/create-service-from-database.go`.

Tous les services sont attachés à `databases_default`. À la création d'une base, ses identifiants sont recopiés dans toutes les applications GitHub existantes, sans association à un projet. Une application compromise dispose ainsi d'accès et de secrets dépassant ses besoins. Aucun plafond mémoire/CPU/PID n'est configuré pour les conteneurs.

**Correction :** réseaux et associations explicites par projet, secrets uniquement pour les consommateurs autorisés, quotas par service. Le modèle actuel doit être considéré comme un environnement de confiance partagé.

### 8. Moyenne — protections de l'API et des opérations longues incomplètes

**Sources :** `internal/api/api.go`, `internal/application/login.go`, `internal/adapter/github.go`, `internal/adapter/docker.go`, `internal/api/handlers/post-deploy.go`.

Pas de limitation de tentatives de connexion ni de taille de corps configurée. Les clients GitHub n'ont pas de délai global, les commandes de build n'ont pas de timeout et les appels Docker utilisent `context.Background()`. Les déploiements sont synchrones, sans verrou par service : deux requêtes peuvent arrêter/supprimer/créer les mêmes conteneurs en parallèle.

**Correction :** limites de débit et de corps, timeouts adaptés, jobs bornés et verrouillés par service, annulation et états persistants des jobs.

### 9. Moyenne — mots de passe générés avec `math/rand`

**Sources :** `internal/utils/random.go`, `internal/utils/env_management.go`.

Le générateur utilisé pour les mots de passe des bases n'est pas cryptographique. Cela ne démontre pas une prédiction triviale de ses sorties, mais il ne doit pas servir à générer des secrets.

**Correction :** utiliser `crypto/rand`, avec tests des contraintes de longueur et alphabet.

## Bugs de déploiement et de gestion

### 10. Élevée — Traefik systématiquement considéré absent

**Source :** `internal/adapter/docker.go:174`.

La détection cherche `/treafik`, alors que la création utilise `traefik`. Chaque déploiement passe par l'arrêt et la recréation du proxy, avec interruption potentielle de toutes les applications. L'erreur de démarrage du proxy est ignorée.

**Correction :** utiliser la constante de nom partout, inspecter l'état réel et propager les erreurs.

### 11. Élevée — redéploiement global et mauvais paramètre HTTP

**Sources :** `internal/api/routes.go:44`, `internal/api/handlers/re-deploy-app.go`, `internal/application/re-deploy-application.go`, `internal/application/deploy-application.go`.

La route déclare `:id`, mais le handler lit `service_name` : la cible reçue est vide. Aucun dépôt n'est alors actualisé, mais tous les services sont redéployés. Même pour un webhook correspondant à un dépôt, la fonction finit par déployer toutes les applications et bases. Les changements d'autres applications peuvent ne pas être récupérés.

**Correction :** sélectionner par ID stable, retourner 404 si absent, actualiser et déployer uniquement la cible et les dépendances nécessaires. Vérifier le nom complet du dépôt et sa branche configurée plutôt que supposer `main`.

### 12. Élevée — build échoué = application déjà arrêtée

**Source :** `internal/application/deploy-application.go:32`.

L'ancienne application est arrêtée avant la construction de la nouvelle image. Si la construction échoue, aucun redémarrage ni rollback n'est effectué. Après démarrage, le contrôle de l'état du conteneur est commenté : un processus qui quitte immédiatement peut malgré tout être enregistré `Running`.

**Correction :** construire d'abord, démarrer et vérifier la nouvelle version, basculer le trafic puis retirer l'ancienne ; conserver une possibilité de retour arrière.

### 13. Élevée — configuration JustDeploy hors du volume Docker documenté

**Sources :** `README.md`, `docker/Dockerfile`, `internal/env.go`.

La commande Docker documentée monte `/app/data`, mais le code utilise `$HOME/.config/justdeploy/database.db`. Le Dockerfile ne redéfinit pas HOME vers ce volume. Avec l'image telle qu'écrite, recréer le conteneur d'administration peut donc perdre configuration, intégration GitHub et compte administrateur, même si `/app/data` est conservé. Cela peut rouvrir l'initialisation publique.

**Correction :** répertoire de données explicite et configurable, montage cohérent, migration et test de recréation.

### 14. Moyenne — imports de dépôts/Compose fragiles

**Sources :** `internal/application/create-app-from-repo.go:22`, `:76`, `internal/adapter/filesystem.go`, `internal/adapter/docker.go:64`.

- Un nom de dépôt non vide sans `/` provoque un accès hors limites après configuration GitHub.
- Un service Compose sans `ports` provoque un accès à `Ports[0]`. C'est une configuration normale pour une base interne.
- `compose.yml` est détecté, mais le parseur ouvre uniquement `docker-compose.yml` ; les variantes `.yaml` ne sont pas gérées.
- L'erreur `yaml.Unmarshal` est ignorée ; volumes, dépendances, contexte de build et commande Compose ne sont pas réellement appliqués par l'import.
- Le hostname généré pour la base (`nom-db`) ne correspond pas au nom de conteneur (`nom-id`) ; aucun alias explicite ne rétablit cette correspondance.
- Un `Dockerfile.xxx` est détecté comme valide, mais le build demande toujours `Dockerfile`.
- Les images des services importés ne sont pas explicitement téléchargées avant leur création, contrairement aux modèles de bases ; Docker peut refuser la création si elles ne sont pas présentes.

**Correction :** valider la forme du dépôt, gérer les champs optionnels et les erreurs, définir précisément le sous-ensemble Compose supporté ou utiliser une interprétation Compose complète. Ces panics peuvent interrompre une requête ; ils ne sont pas assimilés ici à un arrêt certain de tout le serveur.

### 15. Moyenne — endpoints renvoyant un succès sans agir

**Sources :** `internal/application/start-application.go`, `stop-application.go`, `remove-application.go`, `edit-deployement.go`, `get-deploy-list.go`, `get-deploy-config.go`.

Le code opérationnel est commenté. Démarrer, arrêter, supprimer ou modifier un déploiement renvoie `nil` sans action ; des lectures renvoient `nil` ou une configuration vide. Les routes restent enregistrées. La suppression de service a, elle, une implémentation distincte.

**Correction :** implémenter les parcours encore utilisés ou retourner explicitement une erreur d'opération non supportée et retirer les fonctions obsolètes de l'interface.

### 16. Moyenne — erreurs Docker et stockage masquées

**Sources :** `internal/adapter/docker.go:206`, `:220`, `:409`, `internal/application/create-service-from-database.go`, `internal/adapter/database/database.go`.

Les opérations Stop/Remove/Start ignorent leurs erreurs ; certains appels SaveService aussi. Le téléchargement d'une image peut signaler une erreur dans son flux, puis retourner un succès. `PullTreafikImage` et `pullImage` ne propagent pas l'échec. Un nom de base inconnu retourne une configuration vide sans erreur. Des états `ready_to_deploy` ou des suppressions en base peuvent donc être enregistrés alors que l'action Docker a échoué.

**Correction :** remonter les erreurs, distinguer les cas idempotents attendus et préserver la cohérence des états.

### 17. Moyenne — build Nixpacks potentiellement bloqué

**Source :** `internal/adapter/docker.go:120`.

`io.MultiReader(stdout, stderr)` lit d'abord stdout jusqu'à sa fermeture, puis stderr. Si le processus remplit stderr en attendant que le parent le lise, les deux peuvent se bloquer. Les erreurs des pipes et du scanner sont ignorées.

**Correction :** drainer les deux flux simultanément, traiter les erreurs et imposer une échéance au job.

### 18. Moyenne — pertes de réglages dans le formulaire

**Source :** `web/src/components/ServiceDeploySettings.tsx:36`.

Un seul timer est partagé entre visibilité, sous-domaine et variables. Deux modifications rapprochées de champs différents annulent la première sauvegarde, et chaque sauvegarde reconstruit l'objet depuis les anciennes props `service`. Le serveur remplace ensuite l'objet complet : des réglages peuvent être perdus malgré leur affichage local. Aucun nettoyage du timer au démontage n'est prévu.

**Correction :** brouillon unique contenant tous les changements, sauvegarde de ce brouillon et gestion des versions/conflits ; annuler le timer au démontage.

### 19. Moyenne — types frontend et contrat API incohérents

**Sources :** diagnostics TypeScript, `web/src/services/api.ts`, `web/src/components/ServiceDeploySettings.tsx:27`, `internal/api/handlers/post-create-repo.go`.

Modules absents (`createServiceApi`, `ServerStatus`, `ServerConfigForm`, `deployProjectApi`), mauvais types de props et confusion entre `InfoDto` et `ServerDto`. Certains composants peuvent être obsolètes et non accessibles depuis le bundle actuel ; le build Vite réussi ne garantit pas ces parcours. Le client attend toujours une erreur JSON `{message}`, alors que certains handlers renvoient `{error}` ou une chaîne : le message utilisateur peut être vide, voire son parsing échouer pour une réponse texte.

**Correction :** retirer le code mort, aligner DTO et réponses, exécuter TypeScript en CI avant le build.

### 20. Moyenne — configuration et modèles de services défectueux

**Sources :** `internal/domain/settings.go:5`, `internal/utils/env_management.go`, `internal/adapter/database/database.go`, `internal/application/deploy-application.go`.

Les tags `ClientSecret` et `WebhookSecret` sont mal fermés, constat confirmé par `go vet`. Le générateur d'environnement remplace toute valeur non secrète par le nom de variable en minuscules, ce qui détruit notamment la valeur `DATABASE_URL` du modèle Umami. Le type `analytics` est dirigé vers le déploiement GitHub par le cas par défaut, ce qui lui applique un traitement de source/build inadapté.

**Correction :** réparer les tags, préserver les valeurs explicitement définies, définir les besoins de chaque type et refuser les types inconnus.

## Risques complémentaires et limites à approfondir

- Les secrets sont stockés en clair dans SQLite. Le dossier est créé en 0755 ; les permissions du fichier et des sauvegardes ne sont pas durcies explicitement. Vérifier les permissions réelles, limiter à l'utilisateur de service et prévoir une stratégie de protection des secrets au repos.
- L'API Traefik est activée avec `--api.insecure=true`. Son port 8080 n'est pas publié sur l'hôte dans ce code : il ne faut pas conclure à un dashboard public sur Internet. Il reste à examiner son accessibilité depuis le réseau Docker partagé et à désactiver cette API si inutile.
- Le socket Docker est monté dans Traefik et dans l'installation conteneurisée de JustDeploy. Le montage `:ro` du socket n'apporte pas un filtrage des méthodes de l'API Docker. Une compromission de ces composants aurait un fort impact sur l'hôte ; envisager un proxy de socket avec permissions limitées pour Traefik. Le service systemd n'indique pas d'utilisateur dédié.
- Le parcours GitHub utilise un `state` constant (`1234`) à la création et ne valide pas de nonce dans la page de retour. Une réponse contenant les secrets de l'app est renvoyée au navigateur puis journalisée dans sa console. Limiter la réponse aux champs publics nécessaires et lier les retours à une session et un nonce uniques. Le scénario précis de substitution d'intégration mérite un test dédié.
- Adaptateur SSH : vérification de clé hôte désactivée, clé absente pouvant être déréférencée, `log.Fatalf` sur clé invalide, commande `cat` construite sans échappement et fichiers écrits en 0644. Aucun parcours SSH exposé par les routes actives n'a été démontré : priorité moindre tant que cette fonctionnalité reste inactive.
- Les clones sont placés dans des chemins temporaires enregistrés durablement et ne sont pas nettoyés par les parcours examinés : accumulation de fichiers et échec ultérieur si le système purge ces dossiers. Utiliser un répertoire de travail maîtrisé et une politique de rétention.
- Les logs Docker sont analysés par lignes en retirant huit octets à chaque ligne. Le protocole Docker est structuré par trames ; les messages multilignes/TTY ne sont pas correctement pris en charge. Utiliser le démultiplexage Docker approprié et borner la quantité de logs retournée.
- Les tags d'images flottants (`latest`, `traefik`) et téléchargements de scripts sans version/digest nuisent à la reproductibilité. Le Dockerfile part de Go 1.23 alors que le module exige Go 1.25 ; un téléchargement automatique de toolchain peut masquer l'écart. Aligner explicitement versions, lockfiles et contrôles CI.

## Ordre de correction proposé

1. **Protéger les données et l'accès** : vérifier les sauvegardes, corriger les volumes et le répertoire de données, authentifier les webhooks, protéger l'initialisation, imposer HTTPS sur l'administration, supprimer les secrets des logs et renouveler ceux exposés.
2. **Rendre le déploiement fiable** : respecter la visibilité, cibler le bon service, corriger Traefik, isoler les jobs concurrents, construire avant d'arrêter, vérifier la santé et prévoir le rollback.
3. **Rétablir les contrats** : corriger Compose, modèles de services, endpoints inactifs, formulaires et TypeScript ; propager toutes les erreurs opérationnelles.
4. **Verrouiller les régressions** : tests d'intégration sur installation vierge, faux webhook, service privé, restauration d'une base après recréation, build échoué, déploiements concurrents et réglages rapides. Ajouter analyse statique et scans de dépendances/images en CI.

Les bons éléments existants sont bcrypt pour les mots de passe administrateur, `crypto/rand` pour le secret JWT et les requêtes SQL paramétrées examinées. Ils ne compensent pas les problèmes de persistance et de déclenchement des déploiements décrits plus haut.


## Correctifs sur `fix/security-and-reliability-audit` — 2 novembre 2026

Le contenu précédent reste le constat initial. Les changements de cette branche traitent les 20 points :

| Constats | Traitement |
| --- | --- |
| 1, 13 | Répertoire persistant cohérent, volumes nommés pour les nouvelles bases, reprise des montages existants et refus de recréer une base non migrable sans risque. |
| 2, 11 | HMAC sur corps brut, déduplication SQLite, dépôt complet et branche, redéploiement ciblé par ID ; aucune opération globale sur dépôt inconnu. |
| 3, 9 | Secrets retirés des URL Git et logs de structures ; contextes filtrés, logs de build bornés/expurgés ; génération cryptographique des mots de passe. |
| 4, 7 | Visibilité appliquée dans Docker/Traefik, groupes de connexion explicites et suppression du partage automatique des secrets ; quotas CPU/mémoire/PID. |
| 5, 6, 8 | Initialisation atomique protégée par jeton, accès HTTP distant refusé, TLS/proxy explicites, sessions de huit heures révocables, stockage par onglet, limites de requêtes et jobs sérialisés/bornés. |
| 10, 12 | Détection fiable du proxy, reconfiguration seulement si nécessaire, construction avant arrêt, vérification des redémarrages/santé et retour au conteneur précédent. |
| 14, 16, 17 | Import Compose limité et validé explicitement, formats reconnus, commandes/volumes conservés, erreurs propagées, flux Nixpacks drainés ensemble et timeout. |
| 15 | Start/stop/suppression implémentés par service, listes réelles ; anciennes opérations non supportées renvoient HTTP 410. |
| 18, 19, 20 | Brouillon et sauvegarde explicite avec révision, contrats API/types corrigés, contrôle TypeScript dans le build ; modèles/tags JSON réparés et anciens secrets lisibles. |

Adaptateurs SSH, permissions des fichiers, timeouts réseau, installation et CI ont également été renforcés. Les dépendances Go/frontend ont été mises à jour ; les SDK Docker/Echo/JWT obsolètes ont été retirés.

Validation exécutée :

- `go test -race ./...` et `go vet ./...` réussis.
- Lint, TypeScript et build frontend réussis.
- `bun audit` : aucun avis dans le lockfile final.
- `govulncheck` : aucune vulnérabilité atteignable signalée ; un avis subsiste au niveau d'un module requis, sans chemin de code importé/appelé concerné.
- Construction locale de l'image Docker réussie, sans identifiants privés.
- Test sur Docker réel réussi : une donnée survit à la recréation ; un conteneur qui plante est refusé et l'ancien est relancé. Le test utilise ses propres ressources temporaires et est disponible dans le dépôt.
- Test de l'image réussit sur installation vierge : jeton d'initialisation, SPA embarquée et compte conservé après recréation avec le même volume.
- Le scan Trivy local n'a pas abouti : le téléchargement du scanner a rencontré un accès Docker au trousseau annulé. Le scan de l'image est configuré en CI ; aucun résultat « image sans vulnérabilité » n'est revendiqué.

Les changements sont locaux, sans déploiement de production. La rotation des anciens secrets, les sauvegardes/restaurations des données de production et la configuration TLS restent des opérations propres à chaque serveur. Les jobs sont sérialisés et bornés mais ne constituent pas une file durable. Le contrôle sans HEALTHCHECK vérifie la stabilité du processus ; il ne démontre pas la disponibilité métier. Les parcours GitHub réels et la délivrance de certificats ACME nécessitent une validation avec les comptes/domaines de l'installation. Voir [SECURITY-UPGRADE.md](SECURITY-UPGRADE.md).
