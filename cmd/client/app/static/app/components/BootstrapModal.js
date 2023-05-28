app.component('Modal', {
    template: /*html*/ `
    <div ref="modal" class="modal fade" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1" aria-labelledby="staticBackdropLabel" aria-hidden="true">
        <div class="modal-dialog modal-dialog-scrollable" style="max-width: 1000px">
            <div class="modal-content">
                <div class="modal-header">
                    <slot name="header">
                        <div class="modal-title inline-flex" style="height: 24px">
                            <slot name="title"></slot>
                        </div>
                    </slot>
                </div>
                <div class="modal-body">
                    <slot name="body"></slot>
                </div>
                <div class="modal-footer">
                    <slot name="footer"></slot>
                </div>
            </div>
        </div>
    </div>
    `,
    emits: [ 'show', 'hide' ],
    data() { return { emitDelay: 100 } },
    mounted() { this.modal = new bootstrap.Modal(this.$refs.modal) },
    methods: {
        show() {
            this.modal.show()
            this.emit('show')
        },
        hide() {
            this.modal.hide()
            this.emit('hide')
        },
        emit(eventName) {
            setTimeout(() => this.$emit(eventName), this.emitDelay)
        }
    }
})