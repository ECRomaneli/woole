app.component('Sidebar', {
    template: /*html*/ `
        <div id="sidebar" class="d-flex flex-column flex-shrink-0">
            <div class="d-flex">
                <div class="d-flex p-3 me-2 card-shadow dashboard-button" :class="{ active: !selectedRecord }" @click="show()">
                    <img class="square-24" src="assets/images/settings.svg" alt="dashboard">
                </div>
                <div class="d-flex" style="width: 100%;">
                    <input class="p-3 card-shadow input-search" v-model="inputSearch" :class="{ active: inputSearch !== '' }" placeholder="Search..." type="search" spellcheck="false">
                </div>
            </div>
            <div id="records" class="card-shadow" :class="{ loading: recordList.length === 0 }">
                <div class="list-group list-group-flush scrollarea">
                    <sidebar-item
                        v-for="record in filteredRecordList"
                        :record="record"
                        :key="record.id"
                        :class="{ active: isSelectedRecord(record) }"
                        @click="show(record)"
                    ></sidebar-item>
                </div>
            </div>
        </div>
    `,
    emits: [ 'itemSelected' ],
    props: { maxRecords: Number },

    data() {
        return {
            recordList: [],
            filteredRecordList: [],
            selectedRecord: null,
            inputSearch: ""
        }
    },
    created() {
        this.$bus.on('init', (recs) => {
            recs.reverse()
            this.recordList = recs
            this.filteredRecordList = this.recordList.slice()
        })

        this.$bus.on('update', (rec) => {
            this.recordList.unshift(rec)
            while (this.recordList.length > this.maxRecords) {
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
        isSelectedRecord(record) {
            return this.selectedRecord && this.selectedRecord.id === record.id
        },

        async show(record) {
            if (record === void 0) {
                if (this.selectedRecord !== null) {
                    this.selectedRecord = null
                    this.$emit('itemSelected', null)
                }
                return
            }

            if (this.isSelectedRecord(record.id)) { return }

            if (!record.isFetched && record.response) {
                let resp = await fetch('/record/' + record.id + '/response/body')
                record.response.body = await resp.json()
                record.isFetched = true
            }

            this.selectedRecord = record
            this.$emit('itemSelected', record)
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

app.component('SidebarItem', {
    template: /*html*/ `
        <div class="list-group-item list-group-item-action py-3 lh-tight">
            <div class="d-flex w-100 align-items-center justify-content-between">
                <div class="mb-1">
                    <div v-if="isReplay()" class="bg-replay-badge replay-badge badge mr-4"><img src="assets/images/play.svg" alt="replay" /></div>
                    <span class="badge mr-4" :class="methodBadge()">{{ request.method }}</span>
                    <span class="badge" :class="statusBadge()">{{ response.code }}</span>
                </div>
                <div>
                    <small style="font-size: 10px; margin-right: 3px; color: #bbb;">77ms /</small>
                    <small>{{ record.elapsed }}ms</small>
                </div>
                
            </div>
            <div class="mb-1 small">
                <span class="request-path">{{ ellipsis(request.path) }}</span>
                <span v-if="request.query !== void 0" class="request-query badge" :title="request.query">?</span>
            </div>
        </div>
    `,
    props: {
        record: Object
    },
    data() {
        return {
            request: this.record.request,
            response: this.record.response,
            maxLength: 30
        }
    },
    beforeMount() {
        this.request.query = this.request.url.split('?')[1]
    },
    methods: {
        methodBadge() {
            return "bg-" + this.request.method.toLowerCase()
        },
        statusBadge() {
            return "bg-status-" + parseInt(this.response.code/100)
        },
        ellipsis(path) {
            let maxLength = this.maxLength
            let hasQuery = this.request.query !== void 0

            if (hasQuery) { maxLength -= 4 }
            let result = path.length < maxLength ? path : '...' + path.substring(path.length - maxLength)
            return result + (hasQuery ? ' ' : '')
        },
        isReplay() {
            return this.record.id.indexOf('R') !== -1
        }
    }
})