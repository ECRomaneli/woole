const app = Vue.createApp({
    inject: ['$woole', '$date'],

    data() { return { sessionDetails: {}, selectedRecord: null } },

    created() {
        this.setupStream()
        this.$bus.on('record.replay', this.sendRecord)
        this.$bus.on('record.new', this.sendRecord)
        this.$bus.on('record.clear', this.clearRecords)
    },

    computed: {
        host() {
            return this.sessionDetails.https ? this.sessionDetails.https : this.sessionDetails.http
        }
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
                this.$woole.decodeBody(record.response)
            }
            
            this.selectedRecord = record
        },

        async sendRecord(record) {
            const fn = (rec) => {
                if (rec.type === 'replay') {
                    this.$bus.off('stream.new-record', fn)
                    this.$refs.sidebar.scrollTop()
                    this.$refs.sidebar.showRecord(rec)
                }
            }

            this.$bus.on('stream.new-record', fn)
            
            if (record.clientId !== void 0) {
                await fetch('/record/' + record.clientId + '/replay').catch(this.catchAll)
            } else {
                this.$woole.encodeBody(record.request)
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

        setupStream() {
            const es = new EventSource('record/stream')
            let TenSecondErrorThreshold = 1

            es.addEventListener('session', (event) => {
                const data = JSON.parse(event.data)
                this.sessionDetails = data
            })

            es.addEventListener('start', (event) => {
                if (event.data) {
                    let recs = JSON.parse(event.data)
                    recs.sort((a, b) => b.clientId - a.clientId).forEach(this.setupRecord)
                    this.$bus.trigger('stream.start', recs)
                }
            })

            es.addEventListener('new-record', (event) => {
                if (event.data) {
                    let rec = JSON.parse(event.data)
                    this.setupRecord(rec)
                    this.$bus.trigger('stream.new-record', rec)
                }
            })

            es.addEventListener('update-record', (event) => {
                if (event.data) {
                    let rec = JSON.parse(event.data)
                    this.$bus.trigger('stream.update-record', rec)
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

        catchAll(err) {
            console.warn("Error caught: " + err)
            return { ok: false }
        },

        setupRecord(rec) {
            rec.request.host = this.host
            rec.origin = this.getRecordOrigin(rec)
            rec.createdAt = this.$date.from(rec.createdAtMillis).format('MMM DD, hh:mm:ss A')
            this.$woole.decodeQueryParams(rec.request)
            this.$woole.decodeBody(rec.request)
        },

        getRecordOrigin(rec) {
            if (!rec || !rec.request || !rec.request.header) { return null }

            if (!rec.origin) {
                rec.origin = rec.request.header['Origin'] || rec.request.header['origin']
            }

            return rec.origin
        }
    }
})