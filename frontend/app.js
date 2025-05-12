const { createApp } = Vue

createApp({
  data() {
    return {
      ips: '',
      port: '',
      results: [],
      loading: false,
      error: null
    }
  },
  methods: {
    async checkPorts() {
      if (!this.ips || !this.port) {
        this.error = 'Please enter both IP addresses and port number'
        return
      }

      this.loading = true
      this.error = null
      this.results = []

      const ipList = this.ips.split(/[,\n]/).map(ip => ip.trim()).filter(ip => ip)
      const portNumber = parseInt(this.port)
      
      if (ipList.length === 0) {
        this.error = 'Please enter at least one valid IP address'
        this.loading = false
        return
      }

      if (isNaN(portNumber) || portNumber < 1 || portNumber > 65535) {
        this.error = 'Please enter a valid port number (1-65535)'
        this.loading = false
        return
      }

      try {
        const response = await fetch('http://192.168.19.91:8080/api/check-ports', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            ips: ipList,
            port: portNumber
          })
        })

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`)
        }

        const data = await response.json()
        this.results = data.results
      } catch (error) {
        console.error('Error:', error)
        this.error = 'Error checking ports. Please try again.'
      } finally {
        this.loading = false
      }
    }
  }
}).mount('#app') 