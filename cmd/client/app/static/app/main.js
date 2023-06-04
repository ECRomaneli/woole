const app = Vue.createApp({
    data() { return { sessionDetails: {}, selectedRecord: null } },

    created() {
        this.setupStream()
        this.$bus.on('record.replay', this.sendRecord)
        this.$bus.on('record.new', this.sendRecord)
        this.$bus.on('record.curl', this.createCurl)
        this.$bus.on('record.clear', this.clearRecords)
    },

    methods: {
        async itemSelected(record) {
            if (record === null || record.isFetched || !record.response) {
                this.selectedRecord = record
                return
            }

            let resp = await fetch('/record/' + record.clientId + '/response/body').catch(this.catchAll)
            if (resp.ok && resp.status === 200) {
                record.response.body = await resp.json()
                record.isFetched = true
                this.decodeBody(record.response)
            }
            
            this.selectedRecord = record
        },

        async sendRecord(record) {
            this.$bus.once('stream.update', (rec) => this.$refs.sidebar.showRecord(rec))
            
            if (record.clientId !== void 0) {
                await fetch('/record/' + record.clientId + '/replay').catch(this.catchAll)
            } else {
                this.encodeBody(record.request)
                await fetch('/record/request', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(record.request)
                }).catch(this.catchAll)
            }
        },

        async clearRecords() {
            let resp = await fetch('/record', { method: 'DELETE' }).catch(this.catchAll)
            if (resp.ok && resp.status === 200) {
                this.$bus.trigger('stream.start', [])
            }
        },

        async createCurl(record) {
            let resp = await fetch('/record/' + record.clientId + '/request/curl').catch(this.catchAll)
            record.request.curl = resp.ok && resp.status === 200 ? 
                await resp.json() : "Failed to retrieve cURL"
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
                    recs.sort((a, b) => b.clientId - a.clientId).forEach((rec) => this.decodeBody(rec.request))
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
        },

        catchAll(err) {
            console.warn("Error caught: " + err)
            return { ok: false }
        }
    }
})