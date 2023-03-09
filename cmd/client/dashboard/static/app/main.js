const app = Vue.createApp({
    data() {
        return {
            recordList: [],
            filteredRecordList: [],
            selectedRecord: null,
            config: {},
            auth: {},
            events: [],
            inputSearch: ""
        }
    },
    created() {
        this.setupStream()
        this.$bus.on('init', (recs) => {
            recs.reverse()
            this.recordList = recs
            this.filteredRecordList = this.recordList.slice()
        })

        this.$bus.on('update', (rec) => {
            this.recordList.unshift(rec)
            while (this.recordList.length > this.sessionDetails.maxRecords) {
                this.recordList.pop()
            }

            if (this.matchRequest(rec)) {
                this.filteredRecordList.unshift(rec)
            }
        })

        this.$bus.on('show', this.show)
    },

    watch: {
        inputSearch: function (val, oldVal) {
            if (val === "") {
                this.filteredRecordList = this.recordList
                return
            }

            if (val.indexOf(oldVal) === -1) {
                this.filteredRecordList = this.recordList.filter(this.matchRequest)
                return
            }

            this.filteredRecordList = this.filteredRecordList.filter(this.matchRequest)
        }
    },

    methods: {
        setupStream() {
            let es = new EventSource('record/stream')
            let TenSecondErrorThreshold = 1

            es.addEventListener('sessionDetails', event => {
                const data = JSON.parse(event.data)
                this.sessionDetails = data
                this.$forceUpdate()
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
        },

        isSelectedRecord(record) {
            return this.selectedRecord && this.selectedRecord.id === record.id
        },

        async show(record) {
            if (this.isSelectedRecord(record.id)) { return }

            if (!record.isFetched && record.response) {
                let resp = await fetch('/record/' + record.id + '/response/body')
                record.response.body = await resp.json()
                record.isFetched = true
            }

            this.selectedRecord = record
        },

        matchRequest(rec) {
            if (this.inputSearch === "") { return true }

            let tokens = this.inputSearch.split(" ")
            
            let recClone = JSON.parse(JSON.stringify(rec))
            recClone.response.body = null
            let recJson = JSON.stringify(recClone)

            // TODO: Search all tokens simultaneously
            return tokens.every(token => recJson.indexOf(token) !== -1)
        }
    }
})