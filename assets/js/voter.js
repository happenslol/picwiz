
const mainImg1 = document.getElementById('main-1')
const mainImg2 = document.getElementById('main-2')

const loadedImages = []
const bufferElements = document.getElementsByClassName('preload')

for (let i = 0; i < bufferElements.length; i++) {
    bufferElements[i].addEventListener('load', ev => {
        loadedImages.push({
            ref: bufferElements[i],
            src: bufferElements[i].src,
        })

        console.log(`image ${i} loaded!`)
    })
}
