const app = { global: { counter: 1 }, nextInt: () => app.global.counter++ }

app.vue = Vue.createApp({
    data() {
        return {
            recordList: [],
            selectedRecord: null,
            config: {},
            auth: {}
        }
    },
    created() { this.setupStream() },
    methods: {
        async retry() {
            await fetch('/record/' + this.selectedRecord.id + '/retry', { headers: { 'Cache-Control': 'no-cache' } })
            this.show(this.recordList[0]);
        },

        setupStream() {
            let es = new EventSource('record/stream');

            es.addEventListener('info', event => {
                const data = JSON.parse(event.data);
                this.info = data
            });

            es.addEventListener('records', event => {
                if (event.data) {
                    let data = JSON.parse(event.data);
                    this.recordList = data.reverse();
                }
            });

            es.addEventListener('record', event => {
                let data = JSON.parse(event.data);
                if (data !== null) {
                    this.recordList.unshift(data);
                    while (this.recordList.length > this.info.maxRecords) {
                        this.recordList.pop();
                    }
                }
            });

            es.onerror = () => console.error("Failed to retrieve data from event stream")
        },

        isSelectedRecord(record) {
            return this.selectedRecord && this.selectedRecord.id === record.id
        },

        async show(record) {
            if (this.isSelectedRecord(record.id)) { return }

            if (!record.isFetched && record.response) {
                let resp = await fetch('/record/' + record.id + '/response/body');
                let data = await resp.json();
                record.response.body = data
                record.isFetched = true
            }
            
            this.selectedRecord = record
        }
    }
})

