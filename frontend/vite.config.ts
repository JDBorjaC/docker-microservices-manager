import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
	plugins: [react()],

	server: {

		//B:Listen to every IP address
		//B:Needed so the container can reach the Host OS (via the configured port)
		host: true,
		hmr: {port:5173}
	}
})
