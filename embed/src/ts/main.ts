// https://www.youtube.com/watch?v=PCWaFLy3VUo , https://codepen.io/bradtraversy/pen/wvaXKoK
import usercard from "template/user-card.htm"

const template = document.createElement('template')
template.innerHTML = usercard


class UserCard extends HTMLElement {
    showInfo: boolean

    constructor() {
        super()

        this.showInfo = true

        this.attachShadow({ mode: 'open' })
        this.shadowRoot?.appendChild(template.content.cloneNode(true))
        this.shadowRoot!.querySelector('h3')!.innerText = this.getAttribute('name') as string
        this.shadowRoot!.querySelector('img')!.src = this.getAttribute('avatar') as string
    }

    toggleInfo() {
        this.showInfo = !this.showInfo
        const info = this.shadowRoot?.querySelector('.info') as HTMLElement
        const toggleBtn = this.shadowRoot?.querySelector('#toggle-info') as HTMLButtonElement

        if (this.showInfo) {
            info.style.display = 'block'
            toggleBtn.innerText = 'Hide Info'
        } else {
            info.style.display = 'none'
            toggleBtn.innerText = 'Show Info'
        }
    }

    connectedCallback() {
        this.shadowRoot?.querySelector('#toggle-info')?.addEventListener("click", () => this.toggleInfo())
    }

    disconnectedCallback() {
        this.shadowRoot?.querySelector('#toggle-info')?.removeEventListener("", () => { })
    }
}

window.customElements.define('user-card', UserCard)

class Hello {
    constructor(public name: string) {
        this.name = name
    }

    greet() { return "Hello, " + this.name }
}

export default Hello