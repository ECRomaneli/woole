const app = Vue.createApp({
    data() { return { sessionDetails: {}, selectedRecord: null } },

    created() { this.setupStream() },

    methods: {
        setupStream() {
            let es = new EventSource('record/stream')
            let TenSecondErrorThreshold = 1

            es.addEventListener('sessionDetails', event => {
                const data = JSON.parse(event.data)
                this.sessionDetails = data
            })

            es.addEventListener('records', event => {
                if (event.data) {
                    let data = JSON.parse(event.data)
                    this.$bus.trigger('init', data)
                }
            })

            es.addEventListener('record', event => {
                let data = JSON.parse(event.data)
                if (data !== null) { this.$bus.trigger('update', data) }
            })

            es.onerror = () => {
                if (TenSecondErrorThreshold > 0) {
                    TenSecondErrorThreshold--
                    setTimeout(() => TenSecondErrorThreshold++, 10000)
                } else {
                    es.close()
                    console.error("Tunnel connection closed")
                }
            }
        }
    }
})