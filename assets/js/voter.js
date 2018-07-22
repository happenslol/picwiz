
const mainImg1 = document.getElementById('main-1')
const mainImg2 = document.getElementById('main-2')

const loadedImages = []
const bufferElements = document.getElementsByClassName('preload')

const waitingForImages = true

for (let i = 0; i < bufferElements.length; i++) {
    bufferElements[i].addEventListener('load', ev => {
        addLoadedImage({
            ref: bufferElements[i],
            src: bufferElements[i].src,
        })
    })
}

function addLoadedImage(loaded) {
    loadedImages.push(loaded)
}
