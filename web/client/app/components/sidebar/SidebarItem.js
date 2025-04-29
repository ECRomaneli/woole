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
