import js from "@eslint/js";
import tseslint from "typescript-eslint";
import globals from "globals";
import hooks from "eslint-plugin-react-hooks";
import refresh from "eslint-plugin-react-refresh";
export default tseslint.config(
 {ignores:["dist/**",".eslintrc.cjs"]},
 js.configs.recommended,
 ...tseslint.configs.recommended,
 {files:["**/*.{ts,tsx}"],languageOptions:{globals:globals.browser},plugins:{"react-hooks":hooks,"react-refresh":refresh},rules:{
  "react-hooks/rules-of-hooks":"error", "react-hooks/exhaustive-deps":"error",
  "react-refresh/only-export-components":["error",{allowConstantExport:true}],
 }},
 {files:["src/components/ui/button.tsx","src/components/ui/badge.tsx"],rules:{"react-refresh/only-export-components":"off"}},
 {files:["*.config.{js,ts}"],languageOptions:{globals:globals.node}}
);
