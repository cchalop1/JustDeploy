import { Plus, X } from "lucide-react";
import { Env } from "@/services/postFormDetails";

type EnvsManagementsProps = {
  envs: Env[];
  setEnvs: (envs: Env[]) => void;
  canEdit?: boolean;
};

export default function EnvsManagements({ envs, setEnvs, canEdit }: EnvsManagementsProps) {
  const addEnv = () => setEnvs([...envs, { name: "", value: "" }]);

  const removeEnv = (idx: number) => setEnvs(envs.filter((_, i) => i !== idx));

  const updateEnv = (idx: number, field: "name" | "value", val: string) => {
    const updated = [...envs];
    updated[idx] = { ...updated[idx], [field]: val };
    setEnvs(updated);
  };

  return (
    <div className="space-y-1.5">
      {envs.length === 0 && (
        <p className="text-xs text-gray-400 py-2">No environment variables defined.</p>
      )}

      {envs.map((env, idx) => (
        <div key={idx} className="flex items-center gap-2">
          <input
            className="flex-1 min-w-0 h-8 px-2.5 rounded-md border border-gray-200 bg-gray-50 font-mono text-xs text-gray-800 placeholder:text-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-300"
            placeholder="KEY"
            readOnly={!canEdit}
            autoComplete="off"
            value={env.name}
            onChange={(e) => updateEnv(idx, "name", e.target.value)}
          />
          <input
            className="flex-1 min-w-0 h-8 px-2.5 rounded-md border border-gray-200 bg-gray-50 font-mono text-xs text-gray-700 placeholder:text-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-300"
            placeholder="value"
            readOnly={!canEdit}
            autoComplete="off"
            value={env.value}
            onChange={(e) => updateEnv(idx, "value", e.target.value)}
          />
          {canEdit && (
            <button
              onClick={() => removeEnv(idx)}
              className="flex-shrink-0 w-6 h-6 flex items-center justify-center rounded text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors"
            >
              <X className="w-3.5 h-3.5" />
            </button>
          )}
        </div>
      ))}

      {canEdit && (
        <button
          onClick={addEnv}
          className="flex items-center gap-1.5 mt-2 text-xs text-gray-400 hover:text-gray-700 transition-colors"
        >
          <Plus className="w-3.5 h-3.5" />
          Add variable
        </button>
      )}
    </div>
  );
}
