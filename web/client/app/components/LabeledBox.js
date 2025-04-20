app.component('Box', {
    template: /*html*/ `
        <div class="pe-2 pb-2">
            <div ref="box" class="box h-100" :class="{ 'maximized': maximized, 'transparent': transparent }">
                <div class="d-flex justify-content-between align-items-center py-2 px-1 border-bottom">
                    <div class="d-inline-flex pe-none">
                        <img v-if="labelImg" class="svg-icon square-24 ms-2" :src="$image.src(labelImg)" :alt="label">
                        <span class="h5 m-0 ms-2">{{ label }}</span>
                    </div>
                    <div class="btn-toolbar">
                        <slot name="buttons"></slot>
                        <div v-if="maximizable === true" class="maximize-btn ms-3 me-2" @click="toggleView()">
                            <img class="svg-icon square-24" :src="$image.src(maximized ? 'minimize' : 'maximize')" alt="toggle view" />
                        </div>
                    </div>
                </div>
                <div class="box-body py-3 px-2">
                    <slot name="body"></slot>
                </div>
            </div>
        </div>
    `,
    inject: ['$image'],
    props: { labelImg: String, label: String, transparent: Boolean, maximizable: { default: true, type: Boolean } },
    data() { return { maximized: false } },

    methods: { 
        toggleView() {
            this.maximized = !this.maximized;

            // Workaround to make ACE Editor re-wrap lines
            setTimeout(() => window.dispatchEvent(new Event('resize')), 10)
        }
    }
})