import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
<<<<<<< HEAD
import tailwindcss from '@tailwindcss/vite';
import path from "path";

export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
  ],
=======

export default defineConfig({
  plugins: [react()],
>>>>>>> 3913f28a646f762fd92ac93f16be93d0ad6d3ceb
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
<<<<<<< HEAD
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
=======
>>>>>>> 3913f28a646f762fd92ac93f16be93d0ad6d3ceb
});
