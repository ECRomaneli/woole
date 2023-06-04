app.component('Box', {
    template: /*html*/ `
            <div ref="box" class="box px-2 pb-5" :class="getClasses()">
                <div class="d-flex justify-content-between align-items-center py-3 mb-3 border-bottom">
                    <div class="d-inline-flex pe-none">
                        <img class="svg-icon square-24 me-2 ms-2" :src="$image.src(labelImg)" :alt="label">
                        <span class="h5 m-0">{{ label }}</span>
                    </div>
                    <div class="btn-toolbar">
                        <slot name="buttons"></slot>
                        <div class="maximize-btn ms-3 me-2" @click="maximize()">
                            <img class="svg-icon square-24" :src="$image.src(mode.img)" :alt="mode.img" />
                        </div>
                    </div>
                </div>
                <slot name="body"></slot>
            </div>
        </div>
    `,
    inject: ['$image'],
    props: { labelImg: String, label: String, class: String },
    data() { return { mode: { img: "maximize", class: "" }  } },

    methods: { 
        maximize() {
            if (!this.mode.class) {
                this.mode.img = "minimize"
                this.mode.class = "maximized"
            } else {
                this.mode.img = "maximize"
                this.mode.class = ""
            }

            // Workaround to make ACE Editor re-wrap lines
            setTimeout(() => window.dispatchEvent(new Event('resize')), 10)
        },

        getClasses() {
            return this.class + " " + this.mode.class
        }
    }
})