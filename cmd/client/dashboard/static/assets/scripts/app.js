const app = { global: { counter: 1 }, nextInt: () => app.global.counter++ }

app.vue = Vue.createApp({
    data() {
        return {
            recordList: [],
            selectedRecord: null,
            config: {},
            auth: {},
            onUpdateList: []
        }
    },
    created() { this.setupStream() },

    methods: {
        async replay(record) {
            this.onceOnUpdate((rec) => this.show(rec));
            await fetch('/record/' + record.id + '/replay');
        },

        setupStream() {
            let es = new EventSource('record/stream');

            es.addEventListener('info', event => {
                const data = JSON.parse(event.data);
                this.info = data;
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
                    this.onUpdateList.forEach(fn => fn(data));
                }
            });

            es.onerror = (err) => console.error(err)
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
        },

        onceOnUpdate(fn) {
            let onceFn = (record) => {
                const index = this.onUpdateList.indexOf(onceFn);
                this.onUpdateList.splice(index, 1);
                fn(record);
            };
            this.onUpdateList.push(onceFn);
        },

        onUpdate(fn) {
            this.onUpdateList.push(fn);
        }
    }
})

