const app = Vue.createApp({
    data() { return { sessionDetails: {}, selectedRecord: null } },

    created() {
        this.setupStream()
        this.$bus.on('record.replay', this.sendRecord)
        this.$bus.on('record.new', this.sendRecord)
    },

    methods: {
        async itemSelected(record) {
            if (record !== null && !record.isFetched && record.response) {
                let resp = await fetch('/record/' + record.id + '/response/body')
                record.response.body = await resp.json()
                record.isFetched = true
                this.decodeBody(record.response)
            }
            this.selectedRecord = record
        },

        sendRecord(record) {
            this.$bus.once('stream.update', (rec) => this.$refs.sidebar.show(rec))
            
            if (record.id !== void 0) {
                fetch('/record/' + record.id + '/replay')
            } else {
                this.encodeBody(record.request)
                fetch('/record/request', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(record.request)
                })
            }
        },

        setupStream() {
            let es = new EventSource('record/stream')
            let TenSecondErrorThreshold = 1

            es.addEventListener('session', (event) => {
                const data = JSON.parse(event.data)
                this.sessionDetails = data
            })

            es.addEventListener('start', (event) => {
                if (event.data) {
                    let recs = JSON.parse(event.data)
                    recs.forEach((rec) => this.decodeBody(rec.request))
                    this.$bus.trigger('stream.start', recs)
                }
            })

            es.addEventListener('update', (event) => {
                if (event.data) {
                    let rec = JSON.parse(event.data)
                    this.decodeBody(rec.request)
                    this.$bus.trigger('stream.update', rec)
                }
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
        },

        decodeBody(item) {
            if (item.body !== void 0) {
                item.b64Body = item.body
                item.body = atob(item.b64Body)
            }
        },

        encodeBody(item) {
            if (item.b64Body !== void 0) {
                item.body = btoa(item.b64Body)
                item.b64Body = void 0
            }
        }
    }
})