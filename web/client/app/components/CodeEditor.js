app.component('CodeEditor', {
    template: /*html*/ `
        <div class='editor-container'>
            <div class="editor-toolbar ps-2" v-if="pretty.enabled">
                <div v-if="readOnly">
                    <button class="fw-light px-2 m-1" :class="{ active: tab === 'raw' }" @click="changeTab('raw')">Raw</button>
                    <button class="fw-light px-2 m-1" :class="{ active: tab === 'pretty' }" @click="changeTab('pretty')">Pretty</button>
                </div>
                <div v-else>
                    <button class="fw-light px-2 m-1" @click="beautify(); setCode(pretty.code, false)">Beautify</button>
                </div>
            </div>
            <div ref="container"></div>
        </div>
    `,
    inject: [ '$beautifier' ],
    props: { type: String, code: String, readOnly: Boolean, minLines: Number, maxLines: Number },
    data() {
        return {
            typesByMode: {
                      'html': [ 'html', 'xml' ],
                       'css': [ 'css', 'sass', 'scss' ],
                'javascript': [ 'javascript', 'js' ],
                     'json5': [ 'json' ],
                        'sh': [ 'shellscript', 'sh' ]
            },
            tab: 'raw',
            pretty: { enabled: this.enablePretty(), code: null },
            themeElement: document.querySelector('[data-theme]')
        }
    },
    mounted() {
        this.createEditor()
        this.$bus.on('theme.change', this.updateTheme)
    },
    beforeUnmount() {
        this.$bus.off('theme.change', this.updateTheme)
        this.editor.destroy()
    },
    watch: {
        type: function (val) { this.setEditorMode(val) },
        code: function (val) { this.setCode(val) }
    },
    methods: {
        createEditor() {
            this.editor = ace.edit(this.$refs.container, {
                useWorker: false,
                readOnly: this.readOnly,
                autoScrollEditorIntoView: true,
                minLines: this.minLines,
                maxLines: this.maxLines,
                //height: '100%',
                wrap: true
            })
            this.updateTheme()
            if (this.readOnly) {
                this.editor.renderer.$cursorLayer.element.style.display = "none"
            }

            this.setEditorMode(this.type)
            this.setCode(this.code)
        },

        setEditorMode(type) {
            this.pretty.enabled = this.enablePretty()
            if (type) {
                for (const mode in this.typesByMode) {
                    if (this.typesByMode[mode].some(t => this.isType(t))) {
                        this.editor.session.setMode('ace/mode/' + mode)
                        return
                    }
                }
            }
            this.editor.session.setMode('')
        },

        getCode() {
            return this.editor.getValue()
        },

        setCode(code, resetPretty) {
            if (resetPretty !== false) {
                this.tab = 'raw'
                this.pretty.code = null
            }
            this.editor.setValue(code)
            this.editor.gotoLine(1)
        },

        getLength() {
            return this.editor.getValue().length
        },

        forceUpdate() {
            this.editor.renderer.updateFull()
        },

        updateTheme() {
            if (this.themeElement.getAttribute('data-theme') === 'dark') {
                this.editor.setTheme('ace/theme/twilight')
            } else {
                this.editor.setTheme('ace/theme/dawn')
            }
        },

        enablePretty() {
            return this.type && this.$beautifier.supports(this.type)
        },

        beautify() {
            this.pretty.code = this.$beautifier.beautify(this.type, this.getCode())
        },

        changeTab(tab) {
            this.tab = tab
            let code = this.code

            if (tab === 'pretty') {
                if (this.pretty.code === null) { this.beautify() }
                code = this.pretty.code
            }
            
            this.editor.setValue(code)
            this.editor.gotoLine(1)
        },

        isType(type) {
            return this.type.indexOf(type) !== -1
        }
    }
})