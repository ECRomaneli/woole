app.component('Sidebar', {
    template: /*html*/ `
        <nav id="sidebar" class="d-flex flex-column my-2 ms-2">
            <div class="d-flex mb-2">
                <div class="d-flex me-2 sidebar-btn w-100" :class="{ active: !selectedRecord }" @click="showRecord()">
                    <img class="svg-icon square-20" :src="$image.src('view-grid')" alt="settings">
                </div>
                <div class="d-flex me-2 sidebar-btn w-100" @click="toggleTheme()">
                    <img class="svg-icon square-20" :src="$image.src(themeImg)" alt="theme">
                </div>
                <div class="d-flex sidebar-btn w-100" @click="$refs.reqEditor.show()">
                    <img class="svg-icon square-20" :src="$image.src('file-signature')" alt="new request">
                </div>
            </div>
            <div class="d-flex mb-2">
                <input class="d-flex me-2 px-3 w-100 input-search" v-model="inputSearch" :class="{ active: inputSearch !== '' }" placeholder="Search..." type="search" spellcheck="false">
                <div class="d-flex sidebar-btn" @click="clearRecords()">
                    <img class="svg-icon square-24" :src="$image.src('trash2')" alt="clear all" title="Clear All">
                </div>
            </div>
            
            <div id="record-list" class="to-be-removed-h-100" :class="{ loading: recordList.length === 0 }">
                <div class="scrollarea">
                    <sidebar-item
                        v-for="record in filteredRecordList"
                        :record="record"
                        :key="record.clientId"
                        :class="{ active: isSelectedRecord(record) }"
                        @click="showRecord(record)"
                    ></sidebar-item>
                </div>
            </div>
            <request-editor ref="reqEditor"></request-editor>
        </nav>
    `,
    inject: [ '$search', '$image' ],
    emits: [ 'itemSelected' ],
    props: { maxRecords: Number },

    data() {
        return {
            recordList: [],
            filteredRecordList: [],
            selectedRecord: null,
            inputSearch: "",
            themeImg: localStorage.getItem('_woole_theme') ?? 'moon',
            appElement: document.getElementById('app'),
            excludeFromSearch: ['b64body', 'response.body']
        }
    },
    beforeMount() { this.setTheme() },
    created() {
        this.$bus.on('stream.start', (recs) => {
            this.recordList = recs
            this.filteredRecordList = this.recordList.slice()
            this.showRecord()
            this.filter(this.recordList)
        })

        this.$bus.on('stream.update', (rec) => {
            this.recordList.unshift(rec)

            if (this.recordList.length <= this.maxRecords) {
                this.filter([rec], true)
                return
            }

            while (this.recordList.length > this.maxRecords) {
                this.recordList.pop()
            }

            this.filter(this.recordList)
        })
    },

    watch: {
        inputSearch() {
            this.filter(this.recordList)
        }
    },

    methods: {
        isSelectedRecord(record) {
            return this.selectedRecord && this.selectedRecord.clientId === record.clientId
        },

        async showRecord(record) {
            if (record === void 0) {
                if (this.selectedRecord !== null) {
                    this.selectedRecord = null
                    this.$emit('itemSelected', null)
                }
                return
            }

            if (this.isSelectedRecord(record.clientId)) { return }

            this.selectedRecord = record
            this.$emit('itemSelected', record)
        },

        clearRecords() {
            this.$bus.trigger('record.clear')
        },

        filter(recordList, append) {
            let filteredRecordList = this.$search(recordList, this.inputSearch, this.excludeFromSearch)
            if (append) {
                filteredRecordList.reverse()
                filteredRecordList.forEach((rec) => this.filteredRecordList.unshift(rec))
            } else {
                this.filteredRecordList = filteredRecordList
            }
        },

        toggleTheme() {
            this.themeImg = this.themeImg === 'sun' ? 'moon' : 'sun'
            this.setTheme()
        },

        setTheme() {
            if (this.themeImg === 'sun') {
                this.appElement.setAttribute('data-theme', 'light')
                localStorage.setItem('_woole_theme', 'sun')
            } else {
                this.appElement.setAttribute('data-theme', 'dark')
                localStorage.setItem('_woole_theme', 'moon')
            }
        }
    }
})

app.component('SidebarItem', {
    template: /*html*/ `
        <div :client-id="record.clientId" class="record-item p-3 lh-sm">
            <div class="d-flex w-100 mb-2 justify-content-between small">
                <div>
                    <div v-if="record.type === 'replay'" class="bg-replay-badge badge me-1"><img src="assets/images/play.svg" alt="replay" /></div>
                    <div v-else-if="record.type === 'redirect'" class="bg-redirect-badge badge me-1"><img src="assets/images/windows.svg" alt="redirect" /></div>
                    <span class="badge me-1" :class="methodBadge()">{{ request.method }}</span>
                    <span class="badge" :class="statusBadge()">{{ response.code }}</span>
                </div>
                <div class="opacity-50">
                    <small class="fw-light">77ms /&nbsp;</small>
                    <small class="fw-bolder">{{ record.response.elapsed }}ms</small>
                </div>
            </div>
            <div class="mb-1 smallest font-monospace text-end">
                <span>{{ ellipsis(request.path) }}</span>
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
        }
    }
})