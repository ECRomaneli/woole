app.vue.component('RecordViewer', {
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
                                    <button type="button" class="btn btn-sm btn-outline-secondary" @click="replay(record)">Replay</button>
                                    <button type="button" class="btn btn-sm btn-outline-secondary" data-bs-toggle="modal" data-bs-target="#request-submitter" @mouseover="requestSubmitterEnabled = true">w/ Changes</button>
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
                            :titleGroup="[record.response.code, httpStatusMessage[record.response.code], record.response.proto]" 
                            :item="record.response">
                        </record-item>
                    </div>
                </div>
            </div>
            <request-submitter v-if="requestSubmitterEnabled" :modalId="'request-submitter'" :originalRequest="record.request"></request-submitter>
        </div>
    `,
    props: {
        record: Object
    },
    data() {
        return {
            httpStatusMessage: {
                200: 'OK', 201: 'Created', 202: 'Accepted', 203: 'Non-Authoritative Information', 204: 'No Content', 205: 'Reset Content', 206: 'Partial Content',
                300: 'Multiple Choices', 301: 'Moved Permanently', 302: 'Found', 303: 'See Other', 304: 'Not Modified', 305: 'Use Proxy', 307: 'Temporary Redirect',
                400: 'Bad Request', 401: 'Unauthorized', 402: 'Payment Required', 403: 'Forbidden', 404: 'Not Found', 405: 'Method Not Allowed', 406: 'Not Acceptable',
                407: 'Proxy Authentication Required', 408: 'Request Timeout', 409: 'Conflict', 410: 'Gone', 411: 'Length Required', 412: 'Precondition Failed',
                413: 'Request Entity Too Large', 414: 'Request-URI Too Long', 415: 'Unsupported Media Type', 416: 'Requested Range Not Satisfiable', 417: 'Expectation Failed',
                500: 'Internal Server Error', 501: 'Not Implemented', 502: 'Bad Gateway', 503: 'Service Unavailable', 504: 'Gateway Timeout', 505: 'HTTP Version Not Supported'
            },
            maximizeSvg: "assets/images/maximize.svg",
            maximized: "",
            requestSubmitterEnabled: false
        }
    },

    watch: {
        record: function () {
            this.requestSubmitterEnabled = false;
        }
    },

    methods: { 
        replay(record) {
            app.once('update', (rec) => this.$parent.show(rec));
            
            if (record.id !== void 0) {
                fetch('/record/' + record.id + '/replay');
            } else {
                fetch('/record/request/new', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(record.request)
                });
            }
        },

        maximize(card) {
            if (this.maximized === "") {
                this.maximized = card;
                this.maximizeSvg = "assets/images/minimize.svg";
            } else {
                this.maximized = "";
                this.maximizeSvg = "assets/images/maximize.svg";
            }
            // Workaround to make ACE Editor re-wrap lines
            setTimeout(() => window.dispatchEvent(new Event('resize')), 10)
        }
    }
})

app.vue.component('RecordItem', {
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
                <li class="nav-item" role="presentation" v-if="hasBody()" @click='tab=1'>
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
                    <content-editor :content="content" :readOnly="true" :minLines="2" :maxLines="40"></content-editor>
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

            let tokens = contentType.Val.join(";").toLowerCase().split(";").map(str => str.trim())

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

app.vue.component('ContentPreview', {
    template: /*html*/ `
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
    template: /*html*/ `<div :id="id"></div>`,
    props: { content: Object, readOnly: Boolean, minLines: Number, maxLines: Number },
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
                theme: "ace/theme/twilight",
                readOnly: this.readOnly,
                autoScrollEditorIntoView: true,
                minLines: this.minLines,
                maxLines: this.maxLines,
                wrap: true
            });
            if (this.readOnly) {
                this.editor.renderer.$cursorLayer.element.style.display = "none"
            }

            this.updateData();
        },

        updateData() {
            if (this.lastValue === this.content.data) {
                // Workaround to update content when tab is shown after content change
                this.editor.renderer.updateFull()
                return
            }

            if (this.content.type) { this.setEditorMode() }

            this.editor.setValue(atob(this.content.data)/*, -1 to scroll top */)
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
        },

        getValue() {
            return btoa(this.editor.getValue());
        }
    }
})

app.vue.component('RequestSubmitter', {
    template: /*html*/ `
    <form @submit.prevent="submit" class="checkout-form">
        <div :id="modalId" class="modal fade" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1" aria-labelledby="staticBackdropLabel" aria-hidden="true">
            <div class="modal-dialog modal-dialog-scrollable" style="max-width: 1000px">
                <div class="modal-content">
                    <div class="modal-header">
                        <div class="modal-title inline-flex" style="height: 24px">
                            <img class="square-24 me-2" src="assets/images/request.svg" alt="request">
                            <span class="h5">Request</span>
                        </div>
                    </div>
                    <div class="modal-body">
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
                                    <td><textarea placeholder="Name" class="auto-resize" spellcheck="false" @focus="autoResize" @input="autoResize" @blur="autoResize" v-model="header.name"></textarea></td>
                                    <td><textarea placeholder="Value" class="auto-resize" spellcheck="false" @focus="autoResize" @input="autoResize" @blur="autoResize" v-model="header.value"></textarea></td>
                                </tr>
                                <tr>
                                    <td colspan='3'><div class="clickable-img" @click="add()"><img class="square-24" src="assets/images/plus.svg" alt="add-header" style="width: 24px"></div></td>
                                </tr>
                            </tbody>
                        </table>
                        <div class="h5 centered-title">Body</div>
                        <content-editor ref="bodyEditor" :content="content" :readOnly="false" :minLines="20" :maxLines="40"></content-editor>
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal" @click="cancel()">Cancel</button>
                        <button type="submit" class="btn btn-secondary" data-bs-dismiss="modal">Submit</button>
                    </div>
                </div>
            </div>
        </div>
    </form>
    `,
    props: { modalId: String, originalRequest: Object },

    data() {
        return {
            httpMethods: ["HEAD", "GET", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"],
            content: { data: this.originalRequest.body }
        }
    },

    beforeMount() { this.resetRequest() },

    methods: {
        cancel() {
            this.resetRequest();
            this.$forceUpdate();
        },

        async submit() {
            let req = this.clone(this.request);
            req.header = {};
            req.body = this.$refs.bodyEditor.getValue();
            this.request.header.forEach(h => req.header[h.name] = {Val: [h.value]});
            this.$parent.replay({request:req});
        },

        resetRequest() {
            this.request = this.clone(this.originalRequest);
            this.request.header = [];

            Object.keys(this.originalRequest.header).forEach(headerName => {
                this.request.header.push({
                    name: headerName,
                    value: this.originalRequest.header[headerName].Val.join(';')
                });
            });
        },

        add() {
            this.request.header.push({ name: '', value: '' });
            this.$forceUpdate();
        },

        remove(index) {
            this.request.header.splice(index, 1);
            this.$forceUpdate();
        },

        autoResize(event) {
            let el = event.currentTarget;
            el.style.height = 'auto';
            el.style.height = event.type !== 'blur' ? (el.scrollHeight)+'px' : '';
        },

        clone(obj) {
            return JSON.parse(JSON.stringify(obj));
        }
    }
})