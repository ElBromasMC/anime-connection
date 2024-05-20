
const template = document.createElement("template")
template.innerHTML = `
    <div class="carousel">
        <div class="slider">
            <slot></slot>
        </div>
        <!-- Controls -->
        <div class="buttons hidden">
            <button class="btn btn-prev hidden" type="button">
                <svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12l4-4m-4 4 4 4"/>
                </svg>
            </button>
            <div class="separator"></div>
            <button class="btn btn-next hidden" type="button">
                <svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 12H5m14 0-4 4m4-4-4-4"/>
                </svg>
            </button>
        </div>
        <!-- Indicator -->
        <div class="indicator hidden">
            <template>
                <button class="tab" type="button"><span></span></button>
            </template>
        </div>
    </div>
    
    <style>
        :host {
            display: block;
        }
        .carousel {
            position: relative;
            height: 100%;
            overflow: hidden;
        }
        .slider {
            display: flex;
            position: relative;
            height: 100%;
            transition: all 0.5s;
            & ::slotted(*) {
                position: relative;
                flex-basis: 100%;
                flex-grow: 0;
                flex-shrink: 0;
                min-width: 0;
            }
        }

        /* Controls style */
        .buttons {
            position: absolute;
            bottom: 1.75rem;
            right: 1.75rem;
            display: flex;
            border-radius: 1rem;
            overflow: hidden;
            & .btn {
                display: flex;
                justify-content: center;
                align-items: center;
                padding: 0;
                border: 0;
                width: 3.5rem;
                height: 3.5rem;
                background-color: #1d2747;
                cursor: pointer;
                z-index: 1;
                & svg {
                    width: 2rem;
                    height: 2rem;
                    font-size: 2rem;
                    color: #f9f9f9;
                }
                &:hover svg {
                    color: #4c6085;
                }
            }
            & .separator {
                width: 0.25rem;
                align-self: stretch;
                background-color: #39a0ed;
            }
        }

        /* Indicator style */
        .indicator {
            display: flex;
            justify-content: center;
            align-items: center;
            gap: 0.5rem;

            position: absolute;
            bottom: 0.75rem;
            left: 0;
            right: 0;

            z-index: 1;
        }
        button.tab {
            display: flex;
            justify-content: center;
            align-items: center;
            flex-grow: 0;
            flex-shrink: 0;

            min-width: 0;
            width: 2rem;
            height: 2rem;

            padding: 0;
            border: 0;
            
            background-color: transparent;
            cursor: pointer;
            & span {
                width: 100%;
                height: 0.18rem;
                background-color: #9ca3af;
            }
            &:hover span,
            &.active span {
                background-color: #ffffff;
            }
        }

        .hidden {
            display: none;
        }
    </style>
`


export default class CarouselMod extends HTMLElement {
    constructor() {
        super()
        this.currentSlide = 0
        this.numberSlides = 0
        this.autoplayInterval = null

        // Create and append contents to the shadow DOM
        this.shadow = this.attachShadow({ mode: "open" })
        this.shadow.appendChild(template.content.cloneNode(true))
    }

    static get observedAttributes() {
        return ["autoplay", "controls", "indicator"]
    }

    connectedCallback() {
        const slot = this.shadow.querySelector("slot")

        const prevButton = this.shadow.querySelector(".btn-prev")
        const nextButton = this.shadow.querySelector(".btn-next")

        const indicator = this.shadow.querySelector(".indicator")
        const tabTemplate = indicator.querySelector("template").content.cloneNode(true)

        slot.addEventListener("slotchange", () => {
            this.numberSlides = Array.from(slot.assignedElements()).length

            let tabs = []
            for (let i = 0; i < this.numberSlides; i++) {
                const tab = tabTemplate.cloneNode(true)
                tabs.push(tab)
            }
            indicator.replaceChildren(...tabs);
            Array.from(indicator.children).forEach((tab, i) => {
                tab.addEventListener("click", () => {
                    this.currentSlide = i
                    this.removeAttribute("autoplay")
                    this.updateDOM()
                })
            })
            this.updateDOM()
        })
        prevButton.addEventListener("click", () => {
            const n = this.numberSlides
            this.currentSlide = ((this.currentSlide - 1) % n + n) % n   // This is only modular arithmetic
            this.removeAttribute("autoplay")
            this.updateDOM()
        })
        nextButton.addEventListener("click", () => {
            const n = this.numberSlides
            this.currentSlide = ((this.currentSlide + 1) % n + n) % n
            this.removeAttribute("autoplay")
            this.updateDOM()
        })
    }

    attributeChangedCallback(name, oldValue, newValue) {
        if (name === "autoplay") {
            if (newValue != null) {
                this.autoplayInterval = setInterval(() => {
                    const n = this.numberSlides
                    this.currentSlide = ((this.currentSlide + 1) % n + n) % n
                    this.updateDOM()
                }, parseInt(newValue) || 3000)
            } else {
                if (this.autoplayInterval !== null) clearInterval(this.autoplayInterval)
                this.autoplayInterval = null
            }
        } else if (name === "controls") {
            const buttonsShow = this.shadow.querySelector(".buttons")
            if (newValue != null) {
                buttonsShow.classList.remove("hidden")
            } else {
                buttonsShow.classList.add("hidden")
            }
        } else if (name === "indicator") {
            const indicator = this.shadow.querySelector(".indicator")
            if (newValue != null) {
                indicator.classList.remove("hidden")
            } else {
                indicator.classList.add("hidden")
            }
        }
    }

    updateDOM() {
        const slider = this.shadow.querySelector(".slider")
        const tabs = Array.from(this.shadow.querySelector(".indicator").children)

        slider.style.transform = `translateX(${-100 * (this.currentSlide)}%)`
        tabs.forEach((tab, i) => {
            if (i === this.currentSlide) {
                tab.classList.add("active")
            } else {
                tab.classList.remove("active")
            }
        })
    }
}
