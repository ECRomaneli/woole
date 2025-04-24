app.component('Sidebar', {
    template: /*html*/ `
        <nav id="sidebar" class="d-flex flex-column m-2">
            <div class="d-flex mb-2">
                <div class="d-flex me-2 sidebar-btn w-100" :class="{ active: !selectedRecord }" @click="showRecord()">
                    <img class="svg-icon square-20" :src="$image.src('view-grid')" alt="settings" title="Settings">
                </div>
                <div class="d-flex me-2 sidebar-btn w-100" @click="toggleTheme()">
                    <img class="svg-icon square-20" :src="$image.src(themeImg)" alt="theme" title="Theme">
                </div>
                <div class="d-flex me-2 sidebar-btn w-100" @click="clearRecords()">
                    <img class="svg-icon square-20" :src="$image.src('trash2')" alt="clear all" title="Clear All">
                </div>
                <div class="d-flex sidebar-btn w-100" @click="$refs.reqEditor.show()">
                    <img class="svg-icon square-20" :src="$image.src('file-signature')" alt="new request" title="New Request">
                </div>
            </div>
            <div class="d-flex mb-2">
                <input class="d-flex px-3 w-100 input-search" v-model="inputSearch" :class="{ active: inputSearch !== '' }" placeholder="Filter records" type="search" spellcheck="false">
            </div>
            
            <div id="record-list" :class="{ loading: recordList.length === 0 }">
                <div ref="scrollarea" class="scrollarea">
                    <template v-for="(record, index) in filteredRecordList">
                        
                        <sidebar-item
                            :record="record"
                            :key="record.clientId"
                            :class="{ active: isSelectedRecord(record), 'first-item': isOtherHost(record, filteredRecordList[index - 1]) }"
                            @click="showRecord(record)"
                        ></sidebar-item>

                        <div v-if="isOtherHost(record, filteredRecordList[index + 1])" class="d-flex p-1 mb-2 origin">
                            <div class="smallest font-monospace text-center">
                                <span>{{ record.request.forwardedTo }}</span>
                            </div>
                        </div>

                    </template>
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
            this.showRecord()
            this.filterRecords(this.recordList)
        })

        this.$bus.on('stream.new-record', (rec) => {
            this.recordList.unshift(rec)

            if (!this.maxRecords || this.recordList.length <= this.maxRecords) {
                this.appendRecords([rec])
                return
            }                
    
            while (this.recordList.length > this.maxRecords) {
                this.recordList.pop()
            }
            
            this.filterRecords(this.recordList)
        })

        this.$bus.on('stream.update-record', (update) => {
            const recordUpdated = this.recordList.some(rec => {
                if (rec.clientId === update.clientId) {
                    rec.step = update.step
                    rec.response.serverElapsed = update.response.serverElapsed
                    return true
                }
            })
            
            if (recordUpdated && this.inputSearch.indexOf('serverElapsed') !== -1) {
                this.filterRecords(this.recordList)
            }
        })

        this.$bus.on('sidebar.search', (search) => this.inputSearch = search)
    },

    watch: { inputSearch() { this.filterRecords(this.recordList) } },

    methods: {
        isOtherHost(record, otherRecord) {
            return otherRecord === void 0 || record.request.forwardedTo !== otherRecord.request.forwardedTo
        },

        isSelectedRecord(record) {
            return this.selectedRecord && this.selectedRecord.clientId === record.clientId
        },

        scrollTop() {
            this.$refs.scrollarea.scrollTo(0, 0)
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

        filterRecords(recordList) {
            this.filteredRecordList = this.$search(recordList, this.inputSearch, this.excludeFromSearch)
            this.$emit('filterRecords', this.filteredRecordList)
        },

        appendRecords(recordList) {
            if (!recordList.length) return;
            
            const newFilteredRecords = this.$search(recordList, this.inputSearch, this.excludeFromSearch)
            if (newFilteredRecords.length) {
                // Insert at the beginning of the list
                this.filteredRecordList.unshift(...newFilteredRecords)
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
            this.$bus.trigger('theme.change')
        }
    }
})

app.component('SidebarItem', {
    template: /*html*/ `
        <button :client-id="record.clientId" v-bind="$attrs" class="record-item p-3 lh-sm" @mouseover="showToggle = true" @mouseleave="showToggle = false">
            <div class="d-flex w-100 mb-1 justify-content-between small">
                <div class="badge-group" dir="rtl">
                    <span class="badge px-1" :class="statusBadge()">{{ response.code }}</span>
                    <span class="badge px-1 me-1" :class="methodBadge()">{{ request.method }}</span>
                    <div v-if="record.type === 'redirect'" class="bg-redirect-badge badge me-1" title="Redirect">
                        <img :src="$image.src('windows')" alt="redirect" />
                    </div>
                    <div v-else-if="record.type === 'replay'" class="bg-replay-badge badge me-1" title="Replay">
                        <img :src="$image.src('play')" alt="replay" />
                    </div>
                </div>
                <div class="opacity-50 ms-1">
                    <img v-show="showToggle" :src="$image.src('change')" class="me-1 toggle-time" alt="toggle" @click="toggleInfo($event)" />
                    <template v-if="showCreatedAt">
                        <small class="fw-light">{{ createdAt[0] + ', ' }}</small>
                        <small class="fw-bolder">{{ createdAt[1] }}</small>
                    </template>
                    <template v-else-if="response.serverElapsed">
                        <small class="fw-light" title="Client Elapsed Time">{{ response.elapsed }}ms /&nbsp;</small>
                        <small class="fw-bolder" title="Server Elapsed Time">{{ response.serverElapsed }}ms</small>
                    </template>
                    <small v-else class="fw-bolder" title="Client Elapsed Time">{{ response.elapsed }}ms</small>
                </div>
            </div>
            <div class="mb-1 smallest font-monospace text-end">
                <span>{{ ellipsis(request.path) }}</span>
                <span v-if="hasQuery" class="badge bg-query" :title="requestQuery">?</span>
            </div>
        </button>
    `,
    inject: [ '$image', '$date' ],
    inheritAttrs: false,
    props: {
        record: Object
    },
    data() {
        return {
            showCreatedAt: true,
            showToggle: false,
            maxLength: 30
        }
    },
    computed: {
        request() { return this.record.request },
        response() { return this.record.response },
        createdAt() { return this.record.createdAt.split(', ') },
        hasQuery() { return this.queryParam !== void 0 },
        requestQuery() { return this.hasQuery ? this.request.url.split('?')[1] : '' },
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

            if (this.hasQuery) { maxLength -= 4 }
            let result = path.length < maxLength ? path : '...' + path.substring(path.length - maxLength)
            return result + (this.hasQuery ? ' ' : '')
        },
        toggleInfo(e) {
            e.stopPropagation()
            this.showCreatedAt = !this.showCreatedAt
        }
    }
})