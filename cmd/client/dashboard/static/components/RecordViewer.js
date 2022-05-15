app.vue.component('RecordViewer', {
    template: `
        <div class="container-fluid">
            <div class="row row-custom">
                <div class="col-md-12 col-custom-6">
                    <div class="card card-shadow">
                        <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
                            <span class="h4">Request</span>
                            <div class="btn-toolbar mb-2 mb-md-0">
                                <div class="btn-group me-2">
                                    <button type="button" class="btn btn-sm btn-outline-secondary">Prettify</button>
                                    <button type="button" class="btn btn-sm btn-outline-secondary">Export</button>
                                </div>
                                <button type="button" class="btn btn-sm btn-outline-secondary dropdown-toggle">
                                    Replay
                                </button>
                            </div>
                        </div>
                        <record-item 
                            :titleGroup="[record.request.method, record.request.url, record.request.proto]" 
                            :item="record.request">
                        </record-item>
                    </div>
                </div>
                <div class="col-md-12 col-custom-6">
                    <div class="card card-shadow">
                        <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
                            <span class="h4">Response</span>
                        </div>
                        <record-item 
                            :titleGroup="[record.response.code, httpStatusMessage[record.response.code], record.response.proto]" 
                            :item="record.response">
                        </record-item>
                    </div>
                </div>
            </div>
        </div>
    `,
    props: { record: Object },
    data() {
        return {
            httpStatusMessage: {
                200: 'OK', 201: 'Created', 202: 'Accepted', 203: 'Non-Authoritative Information', 204: 'No Content', 205: 'Reset Content', 206: 'Partial Content',
                300: 'Multiple Choices', 301: 'Moved Permanently', 302: 'Found', 303: 'See Other', 304: 'Not Modified', 305: 'Use Proxy', 307: 'Temporary Redirect',
                400: 'Bad Request', 401: 'Unauthorized', 402: 'Payment Required', 403: 'Forbidden', 404: 'Not Found', 405: 'Method Not Allowed', 406: 'Not Acceptable',
                407: 'Proxy Authentication Required', 408: 'Request Timeout', 409: 'Conflict', 410: 'Gone', 411: 'Length Required', 412: 'Precondition Failed',
                413: 'Request Entity Too Large', 414: 'Request-URI Too Long', 415: 'Unsupported Media Type', 416: 'Requested Range Not Satisfiable', 417: 'Expectation Failed',
                500: 'Internal Server Error', 501: 'Not Implemented', 502: 'Bad Gateway', 503: 'Service Unavailable', 504: 'Gateway Timeout', 505: 'HTTP Version Not Supported'
            }
        }
    }
})

app.vue.component('RecordItem', {
    template: `
        <div class="col-md-12 mb-5">
            <div class="input-group mb-3">
                <span class="input-group-text">{{ titleGroup[0] }}</span>
                <input type="text" class="form-control" disabled :value="titleGroup[1]">
                <span class="input-group-text">{{ titleGroup[2] }}</span>
            </div>

            <ul class="nav nav-tabs" role="tablist">
                <li class="nav-item" role="presentation" @click='tab=0'>
                    <button class="nav-link" :class="{ active: tab === 0 }">Header</button>
                </li>
                <li class="nav-item" role="presentation" v-if="hasBody()" @click="tab=1; bodyFlag=bodyFlag?0:1">
                    <button class="nav-link" :class="{ active: tab === 1 }">Body</button>
                </li>

                <li class="nav-item" role="presentation" v-if="isPreviewSupported()" @click="tab=2">
                    <button class="nav-link" :class="{ active: tab === 2 }">Preview</button>
                </li>
            </ul>

            <div class="tab-content">
                <div class="tab-pane show" :class="{ active: tab === 0 }">
                    <header-grid :header="item.header"></header-grid>
                </div>

                <div class="tab-pane mt-3" :class="{ active: tab === 1 }">
                    <content-editor :content="content" :readOnly="true" :updated="bodyFlag"></content-editor>
                </div>

                <div class="tab-pane mt-3" :class="{ active: tab === 2 }">
                    <content-preview :content="content"></content-preview>
                </div>
            </div>
        </div>
    `,
    props: { titleGroup: Array, item: Object },
    data() { 
        return {
            supportedPreviews: ['video', 'image'],
            tab: 0,
            bodyFlag: 0
        }
    },
    beforeMount() { this.parseBody() },
    beforeUpdate() {
        this.parseBody()
        this.closeUnavailableTab()
    },
    methods: { 
        closeUnavailableTab() {
            let tab = this.tab
            if (tab === 2 && !this.isPreviewSupported()) { tab = 1 }
            if (tab === 1 && !this.hasBody()) { tab = 0 }
            this.tab = tab
        },

        parseBody() {
            if (this.item.body === void 0 || this.item.body === null) { this.item.body = '' }

            this.content = { data: this.item.body, category: '', type: '' }

            if (!this.item || !this.item.header) { return }

            let contentType = this.item.header['Content-Type']
            if (contentType === void 0 || contentType === '' || contentType.length === 0) { return }

            let tokens = contentType.join(";").toLowerCase().split(";").map(str => str.trim())

            // Parse the xxxx/yyyyy content-type
            let categoryAndType = tokens.shift().split('/')
            this.content.category = categoryAndType[0]
            this.content.type = categoryAndType[1]

            // Parse other possible tokens
            for (let token in tokens) {
                token.indexOf("charset=") === 0 && (this.content.charset = token.substring(8))
            }
        },

        hasBody() { return this.content.data !== '' },
        isPreviewSupported() { return this.supportedPreviews.indexOf(this.content.category) !== -1 }
    }
    
})

app.vue.component('HeaderGrid', {
    template: `
        <table class="table table-striped table-hover header-grid" aria-label="headers">
            <thead>
                <tr>
                <th scope="name">Name</th>
                <th scope="value">Value</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(value, name) in header">
                    <th scope="name">{{ name }}</th>
                    <td>{{ value.join(";") }}</td>
                </tr>
            </tbody>
        </table>
    `,
    props: { header: Object }
})

app.vue.component('ContentPreview', {
    template: `
        <div class='content-preview'>
            <img v-if="content.category === 'image'" :src="contentToSource()" alt="preview" />
            <video v-else-if="content.category === 'video'" controls=""><source :type="contentType()" :src="contentToSource()"></video>
        </div>
    `,
    props: { content: Object },
    methods: {
        contentType() {
            return this.content.category + '/' + this.content.type
        },
        // The body is already in base64
        contentToSource() {
            return "data:" + this.contentType() + ";base64," + this.content.data
        }
    }
})

app.vue.component('ContentEditor', {
    template: `<div :id="id"></div>`,
    props: { content: Object, readOnly: Boolean },
    data() {
        return {
            id: "ace-editor-" + app.nextInt(),
            typesByMode: {
                'html': [ 'html', 'xml' ],
                'css': [ 'css', 'sass', 'scss' ],
                'javascript': [ 'javascript' ],
                'json5': [ 'json' ]
            }
        }
    },
    mounted() { this.createEditor() },
    beforeUpdate() { this.updateData() },
    beforeUnmount() { this.editor.destroy() },
    methods: {
        createEditor() {
            this.editor = ace.edit(this.id, {
                useWorker: false,
                theme: "ace/theme/chrome",
                readOnly: this.readOnly,
                autoScrollEditorIntoView: true,
                minLines: 2,
                maxLines: 40,
                wrap: true
            });
            if (this.readOnly) {
                this.editor.renderer.$cursorLayer.element.style.display = "none"
            }
        },

        updateData() {
            if (this.lastValue === this.content.data) {
                this.editor.renderer.updateFull()
                return
            }

            if (this.content) { this.setEditorMode() }

            this.editor.setValue(atob(this.content.data))
            this.editor.clearSelection()

            this.lastValue = this.content.data
        },

        setEditorMode() {
            for (const mode in this.typesByMode) {
                if (this.typesByMode[mode].some(t => this.content.type.indexOf(t) !== -1)) {
                    this.editor.session.setMode('ace/mode/' + mode)
                    return
                }                
            }
            this.editor.session.setMode('')
        }
    }
})