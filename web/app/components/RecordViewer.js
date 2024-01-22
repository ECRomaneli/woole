app.component('RecordViewer', {
    template: /*html*/ `
        <box label-img="request" label="Request" class="w-100">
            <template #buttons>
                <button type="button" class="btn btn-sm" @click="openCurlViewer()">cURL</button>
                <div class="btn-group ms-2 me-2">
                    <button type="button" class="btn btn-sm" @click="$bus.trigger('record.replay', record)">Replay</button>
                    <button type="button" class="btn btn-sm" @click="$refs.reqEditor.show()">w/ Changes</button>
                </div>
                <a v-if="record.request.method.toLowerCase() === 'get'" class="btn me-2 lh-1" :href="getFullUrl()" target="_blank" title="Open in a new tab">
                    <img class="svg-icon square-16 h-100" :src="$image.src('windows2')" alt="redirect">
                </a>
            </template>
            <template #body>
                <record-item 
                    :titleGroup="[record.request.method, record.request.url, record.request.proto, getFullUrl()]" 
                    :item="record.request">
                </record-item>
            </template>
        </box>
        <box label-img="response" label="Response" class="w-100">
            <template #body>
                <record-item 
                    :titleGroup="[record.response.code, $constants.HTTP_STATUS_MESSAGE[record.response.code], record.response.proto]" 
                    :item="record.response">
                </record-item>
            </template>
        </box>

        <request-editor ref="reqEditor" :originalRequest="record.request"></request-editor>
        <code-viewer ref='curlViewer' type="shellscript" :code="record.request.curl"></code-viewer>
    `,
    inject: ['$constants', '$woole', '$image'],
    props: { record: Object },

    methods: {
        openCurlViewer() {
            if (!this.record.request.curl) {
                this.$woole.parseRequestToCurl(this.record.request)
            }
            this.$refs.curlViewer.show()
        },

        getFullUrl() { return this.record.request.host + this.record.request.url }
    }
})

app.component('RecordItem', {
    template: /*html*/ `
        <div class="highlighted-group input-group mb-3" @mouseover="enableCopy = true" @mouseleave="enableCopy = false">
            <span class="input-group-text">{{ titleGroup[0] }}</span>
            <input type="text" class="form-control" disabled :value="titleGroup[1]">
            
            <button v-if="titleGroup[3]" class="btn img-btn" @click="$clipboard.writeText(titleGroup[3])" title="Copy">
                <Transition name="fast-fade">
                    <img v-show="enableCopy" class="svg-icon square-16 ms-2 me-2" :src="$image.src('copy')" alt="copy">
                </Transition>
            </button>
            
            <span v-if="titleGroup[2]" class="input-group-text">{{ titleGroup[2] }}</span>
        </div>

        <ul class="inline-tabs">
            <li v-if="hasHeader()" @click="tab = 'header'">
                <button class="tab" :class="{ active: tab === 'header' }">Header</button>
            </li>
            <li v-if="hasParam()" @click="tab = 'param'">
                <button class="tab" :class="{ active: tab === 'param' }">Params</button>
            </li>
            <li v-if="hasBody()" @click="tab = 'body'; $refs.codeEditor.forceUpdate()">
                <button class="tab" :class="{ active: tab === 'body' }">
                    Body <span v-if="hasBody()" class="badge fw-light" style="font-size:.6rem">{{ bodySize }}</span>
                </button>
            </li>

            <li v-if="isPreviewSupported()" @click="tab = 'preview'">
                <button class="tab" :class="{ active: tab === 'preview' }">Preview</button>
            </li>
        </ul>

        <div class="tab-content">
            <div class="tab-pane show" :class="{ active: tab === 'header' }">
                <map-table :map="item.header"></map-table>
            </div>

            <div class="tab-pane" :class="{ active: tab === 'param' }">
                <map-table :map="item.queryParams"></map-table>
            </div>

            <div class="tab-pane mt-3" :class="{ active: tab === 'body' }">
                <code-editor ref="codeEditor" :type="content.type" :code="item.body" :readOnly="true" :minLines="2" :maxLines="39"></code-editor>
            </div>

            <div class="tab-pane mt-3" :class="{ active: tab === 'preview' }">
                <base64-viewer :category="content.category" :type="content.type" :data="item.b64Body"></base64-viewer>
            </div>
        </div>
    `,
    inject: ['$woole', '$image', '$clipboard'],
    props: { titleGroup: Array, item: Object },
    data() { return { supportedPreviews: ['image', 'video', 'audio'], tab: 'header', enableCopy: false } },
    beforeMount() { this.parseBody() },
    mounted() { this.selectAvailableTab() },
    beforeUpdate() {
        this.parseBody()
        this.selectAvailableTab()
    },
    computed: {
        bodySize() {
            let length = parseInt(this.item.header['Content-Length']) || this.body.length

            if (!length) { return '0 B' }

            const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
            let index = 0;
        
            while (length >= 1000 && index < units.length - 1) {
                length /= 1024;
                index++;
            }
        
            return length.toFixed(2) + ' ' + units[index];
        }
    },
    methods: { 
        selectAvailableTab() {
            let availableTabs = []
            this.isPreviewSupported()   && availableTabs.push('preview')
            this.hasBody()              && availableTabs.push('body')
            this.hasParam()             && availableTabs.push('param')
            this.hasHeader()            && availableTabs.push('header')

            if (availableTabs.some(tab => this.tab === tab)) { return }
            this.tab = availableTabs.length ? availableTabs[0] : ''
        },

        parseBody() {
            this.content = {}
            if (this.item.body === void 0 || this.item.body === null) { this.item.body = '' }
            if (!this.item.header) { return }
            const contentType = this.item.header['Content-Type']
            if (contentType) {
                this.content = this.$woole.parseContentType(contentType)
            }
        },

        async copyUrl() {
            await navigator.clipboard.writeText(this.titleGroup[1])
        },

        hasHeader() { return this.item.header && Object.keys(this.item.header).length > 0 },
        hasParam() { return this.item.queryParams && Object.keys(this.item.queryParams).length > 0 },
        hasBody() { return this.item.body !== '' },
        isPreviewSupported() {
            return this.supportedPreviews.some(c => c === this.content.category) && this.hasBody()
        }
    }
    
})