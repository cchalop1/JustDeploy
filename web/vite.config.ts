import path from "path";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";
import { reactClickToComponent } from "vite-plugin-react-click-to-component";
import Inspect from "vite-plugin-inspect";

export default defineConfig({
  plugins: [react(), reactClickToComponent(), Inspect()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
});
