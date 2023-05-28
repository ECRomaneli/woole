app.component('RecordViewer', {
    template: /*html*/ `
        <div class="container-fluid">
            <div class="row row-custom">
                <div class="col-md-12 col-custom-6">
                    <div class="card card-shadow" :class="{ maximized: maximized === 'request' }">
                        <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
                            <div class="inline-flex">
                                <img class="square-24 me-2 ms-2" src="assets/images/request.svg" alt="request">
                                <span class="h5">Request</span>
                            </div>
                            <div class="btn-toolbar">
                                <div class="btn-group">
                                    <button type="button" class="btn btn-sm btn-outline-secondary" @click="openCurlViewer()">cURL</button>
                                </div>
                                <div class="btn-group ms-2 me-2">
                                    <button type="button" class="btn btn-sm btn-outline-secondary" @click="$bus.trigger('record.replay', record)">Replay</button>
                                    <button type="button" class="btn btn-sm btn-outline-secondary" @click="openRequestEditor()">w/ Changes</button>
                                </div>
                                <div class="maximize-btn ms-3 me-2" @click="maximize('request')">
                                    <img class="square-24" :src="maximizeSvg" alt="maximize" />
                                </div>
                            </div>
                        </div>
                        <record-item 
                            :titleGroup="[record.request.method, record.request.url, record.request.proto]" 
                            :item="record.request">
                        </record-item>
                    </div>
                </div>
                <div class="col-md-12 col-custom-6">
                    <div class="card card-shadow" :class="{ maximized: maximized === 'response' }">
                        <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
                            <div class="inline-flex">
                                <img class="square-24 me-2 ms-2" src="assets/images/response.svg" alt="response">
                                <span class="h5">Response</span>
                            </div>
                            <div class="btn-toolbar">
                                <div class="maximize-btn ms-3 me-2" @click="maximize('response')">
                                    <img class="square-24" :src="maximizeSvg" alt="maximize" />
                                </div>
                            </div>
                        </div>
                        <record-item 
                            :titleGroup="[record.response.code, $constants.HTTP_STATUS_MESSAGE[record.response.code], record.response.proto]" 
                            :item="record.response">
                        </record-item>
                    </div>
                </div>
            </div>
            <request-editor ref="reqEditor" :originalRequest="record.request"></request-editor>
            <code-viewer ref='curlViewer' :type="'shellscript'" :code="requestCurl"></code-viewer>
        </div>
    `,
    inject: ['$constants'],
    props: { record: Object },
    data() {
        return {
            maximizeSvg: "assets/images/maximize.svg",
            maximized: "",
            requestCurl: ''
        }
    },

    watch: {
        record: function () {
            this.requestCurl = ''
        }
    },

    methods: { 
        maximize(card) {
            if (this.maximized === "") {
                this.maximized = card
                this.maximizeSvg = "assets/images/minimize.svg"
            } else {
                this.maximized = ""
                this.maximizeSvg = "assets/images/maximize.svg"
            }
            // Workaround to make ACE Editor re-wrap lines
            setTimeout(() => window.dispatchEvent(new Event('resize')), 10)
        },

        openRequestEditor() {
            this.$refs.reqEditor.show()
        },

        openCurlViewer() {
            if (!this.requestCurl) {
                this.$bus.once('record.curl.retrieved', curl => this.requestCurl = curl)
                this.$bus.trigger('record.curl', this.record)
            }
            this.$refs.curlViewer.show()
        }
    }
})

app.component('RecordItem', {
    template: /*html*/ `
        <div class="col-md-12 mb-5">
            <div class="cursor-default input-group mb-3">
                <span class="input-group-text">{{ titleGroup[0] }}</span>
                <input type="text" class="form-control" disabled :value="titleGroup[1]">
                <span class="input-group-text">{{ titleGroup[2] }}</span>
            </div>

            <ul class="nav nav-tabs" role="tablist">
                <li class="nav-item" role="presentation" @click='tab=0'>
                    <button class="nav-link" :class="{ active: tab === 0 }">Header</button>
                </li>
                <li class="nav-item" role="presentation" v-if="hasBody()" @click='tab=1; $refs.codeEditor.forceUpdate()'>
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
                    <code-editor ref="codeEditor" :type="content.type" :code="item.body" :readOnly="true" :minLines="2" :maxLines="40"></code-editor>
                </div>

                <div class="tab-pane mt-3" :class="{ active: tab === 2 }">
                    <base64-viewer :category="content.category" :type="content.type" :data="item.b64Body"></base64-viewer>
                </div>
            </div>
        </div>
    `,
    inject: ['$util'],
    props: { titleGroup: Array, item: Object },
    data() { return { supportedPreviews: ['video', 'image'], tab: 0 } },
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
            this.content = {}
            if (this.item.body === void 0 || this.item.body === null) { this.item.body = '' }
            if (!this.item.header) { return }
            let contentType = this.item.header['Content-Type']
            if (!contentType) { return }
            this.content = this.$util.parseContentType(contentType.Val.join(";"))
        },

        hasBody() { return this.item.body !== '' },
        isPreviewSupported() { return this.supportedPreviews.indexOf(this.content.category) !== -1 }
    }
    
})

app.component('HeaderGrid', {
    template: /*html*/ `
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
                    <td>{{ value.Val.join(";") }}</td>
                </tr>
            </tbody>
        </table>
    `,
    props: { header: Object }
})

app.component('RequestEditor', {
    template: /*html*/ `
        <modal ref="modal" @show="cloneRequest()" @hide="removeRequest()">
            <template #title>
                <img class="square-24 me-2" src="assets/images/request.svg" alt="request">
                <span class="h5">Request</span>
            </template>
            <template #body v-if="request">
                <div class="input-group mb-3">
                    <select class="request-method input-group-text" name="method" v-model="request.method">
                        <option v-for="(method) in httpMethods" :value="method">{{ method }}</option>
                    </select>
                    <input name="url" type="text" class="form-control" spellcheck="false" aria-label="url" v-model="request.url">
                </div>
                <div class="h5 centered-title">Headers</div>
                <table class="table table-striped table-hover header-grid" aria-label="headers">
                    <thead>
                        <tr><th scope="remove"></th><th scope="name">Name</th><th scope="value">Value</th></tr>
                    </thead>
                    <tbody>
                        <tr v-for="(header, index) in request.header" :key="index">
                            <td><div class="clickable-img" @click="remove(index)"><img class="square-24" src="assets/images/trash.svg" alt="remove-header"></div></td>
                            <td><textarea placeholder="Name" class="auto-resize" spellcheck="false" @focus="autoResize" @input="autoResize" @blur="onBlur($event, header)" v-model="header.name"></textarea></td>
                            <td><textarea placeholder="Value" class="auto-resize" spellcheck="false" @focus="autoResize" @input="autoResize" @blur="onBlur($event, header)" v-model="header.value"></textarea></td>
                        </tr>
                        <tr>
                            <td colspan='3'><div class="clickable-img" @click="add()"><img class="square-24" src="assets/images/plus.svg" alt="add-header" style="width: 24px"></div></td>
                        </tr>
                    </tbody>
                </table>
                <div style="padding-left: 8px"><input type="checkbox" v-model="isAutoContentLengthEnabled"><label class="small-label">Calculate "Content-Length" on submit</label></div>
                <div class="h5 centered-title">Body</div>
                <code-editor ref="codeEditor" :code="originalRequest.body" :type="content.type" :readOnly="false" :minLines="20" :maxLines="40"></code-editor>
            </template>
            <template #footer v-if="request">
                <button type="button" class="btn btn-secondary" data-bs-dismiss="modal" @click="$refs.modal.hide()">Cancel</button>
                <button type="button" class="btn btn-secondary" @click="submit()">Submit</button>
            </template>
        </modal>
    `,
    inject: [ '$util', '$clone' ],
    props: { originalRequest: Object },

    data() {
        return {
            httpMethods: ["HEAD", "GET", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"],
            content: {},
            request: null,
            isAutoContentLengthEnabled: true
        }
    },

    methods: {
        async submit() {
            let req = this.$clone(this.request)
            req.header = {}
            req.b64Body = this.$refs.codeEditor.getCode()
            req.path = new URL(req.url, "http://dummy").pathname
            this.request.header.forEach(h => req.header[h.name] = { Val: [h.value] })

            if (this.isAutoContentLengthEnabled) {
                req.header['Content-Length'] = { Val: [`${this.$refs.codeEditor.getLength()}`] }
            }

            this.$bus.trigger('record.new', { request: req })
            this.$refs.modal.hide()
        },

        cloneRequest() {
            this.request = this.$clone(this.originalRequest)
            this.request.header = []

            Object.keys(this.originalRequest.header).forEach(headerName => {
                let newHeader = {
                    name: headerName,
                    value: this.originalRequest.header[headerName].Val.join(';')
                }

                this.request.header.push(newHeader)

                if (newHeader.name.toLowerCase() === 'content-type') {
                    this.content = this.$util.parseContentType(newHeader.value)
                }
            })
        },

        removeRequest() { this.request = null },

        add() {
            this.request.header.push({ name: '', value: '' })
        },

        remove(index) {
            let header = this.request.header.splice(index, 1)[0]
            if (header.name.toLowerCase() === 'content-type') {
                this.content = {}
            }
        },

        autoResize(event) {
            let el = event.currentTarget
            el.style.height = 'auto'
            el.style.height = event.type !== 'blur' ? el.scrollHeight + 'px' : ''
        },

        onBlur(event, header) {
            this.autoResize(event)
            if (header.name.toLowerCase() === 'content-type') {
                this.content = this.$util.parseContentType(header.value)
            }
        },

        show() { this.$refs.modal.show() }
    }
})

app.component('CodeViewer', {
    template: /*html*/ `
    <modal ref="modal">
        <template #title>
            <img class="square-24 me-2" src="assets/images/request.svg" alt="request">
            <span class="h5">Code Viewer</span>
        </template>
        <template #body v-if="code">
            <code-editor ref="editor" :type="type" :code="code" :readOnly="false" :minLines="10" :maxLines="40"></code-editor>
        </template>
        <template #footer v-if="code">
            <button v-if="code" type="button" class="btn btn-secondary" @click="copyToClipboard()">{{ copyBtnText }}</button>
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
        </template>
    </modal>
    `,
    props: { type: String, code: String },
    data() { return { copyBtnText: "Copy to Clipboard" } },
    methods: {
        async copyToClipboard() {
            await navigator.clipboard.writeText(this.$refs.editor.getCode())
            this.copyBtnText = "Copied!"
            setTimeout(() => this.copyBtnText = "Copy to Clipboard", 3000)
        },

        show() { this.$refs.modal.show() }
    }
})