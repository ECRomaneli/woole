const app = { global: { counter: 1 }, nextInt: () => app.global.counter++ }

app.vue = Vue.createApp({
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
        this.setupStream();
        this.once('init', (recs) => {
            recs.reverse();
            this.recordList = recs;
            this.filteredRecordList = this.recordList.slice();
        });

        this.on('update', rec => {
            this.recordList.unshift(rec);
            while (this.recordList.length > this.info.maxRecords) {
                this.recordList.pop();
            }

            if (this.matchRequest(rec)) {
                this.filteredRecordList.unshift(rec);
            }
        });
    },

    watch: {
        inputSearch: function (val, oldVal) {
            if (val === "") {
                this.filteredRecordList = this.recordList;
                return;
            }

            if (val.indexOf(oldVal) === -1) {
                this.filteredRecordList = this.recordList.filter(this.matchRequest);
                return;
            }

            this.filteredRecordList = this.filteredRecordList.filter(this.matchRequest);
        }
    },

    methods: {
        setupStream() {
            let es = new EventSource('record/stream');

            es.addEventListener('info', event => {
                const data = JSON.parse(event.data);
                this.info = data;
                this.$forceUpdate();
            });

            es.addEventListener('records', event => {
                if (event.data) {
                    let data = JSON.parse(event.data);
                    this.trigger('init', data);
                }
            });

            es.addEventListener('record', event => {
                let data = JSON.parse(event.data);
                if (data !== null) { this.trigger('update', data); }
            });

            es.onerror = (err) => console.error(err);
        },

        async replay(record) {
            this.once('update', (rec) => this.show(rec));
            await fetch('/record/' + record.id + '/replay');
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

        matchRequest(rec) {
            if (this.inputSearch === "") { return true; }

            let tokens = this.inputSearch.split(" ");
            
            let recClone = JSON.parse(JSON.stringify(rec));
            recClone.response.body = null;
            let recJson = JSON.stringify(recClone);

            // TODO: Search all tokens simultaneously
            return tokens.every(token => recJson.indexOf(token) !== -1);
        },

        trigger(eventName, data) {
            if (this.events[eventName] !== void 0) {
                this.events[eventName].forEach(fn => fn(data));
            }
        },

        once(eventName, fn) {
            let onceFn = (record) => {
                const index = this.events[eventName].indexOf(onceFn);
                this.events[eventName].splice(index, 1);
                fn(record);
            };
            this.on(eventName, onceFn);
        },

        on(eventName, fn) {
            if (this.events[eventName] === void 0) {
                this.events[eventName] = [];
            }

            this.events[eventName].push(fn);
        }
    }
})

