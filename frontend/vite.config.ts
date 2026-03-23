import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
	plugins: [react()],

	server: {

		//B:Listen to every IP address
		//B:Needed so the container can reach the Host OS (via the configured port)
		host: true,
		hmr: {port:5173},

		//C: Due to WSL2 limitations, polling needs to be enabled so that hmr works. (causes high CPU usage)
		watch: {
			usePolling: true
		}
	},


	
})
