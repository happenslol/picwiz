const mainImg1 = document.getElementById('main-1')
const mainImg2 = document.getElementById('main-2')

let activeBuffer = mainImg1
let backBuffer = mainImg2

let loadedImages = []
let initialized = false
const bufferElements = document.getElementsByClassName('preload')

let waitingForImages = true

window.addEventListener('keyup', ev => {
    if (ev.keyCode === 65) {
        console.log('upvoted')
        voteCurrentPic(true)
    }

    if (ev.keyCode === 83) {
        console.log('downvoted')
        voteCurrentPic(false)
    }
})

for (let i = 0; i < bufferElements.length; i++) {
    bufferElements[i].addEventListener('load', ev => {
        addLoadedImage(bufferElements[i])
    })
}

function addLoadedImage(loaded) {
    loadedImages.push(loaded)

    if (!initialized) {
        if (loadedImages.length < 3) return

        console.log('got first pictures')
        console.dir(loadedImages)
        initialized = true

        const firstPics = loadedImages.splice(0, 2)
        activeBuffer.src = firstPics[0].src
        backBuffer.src = firstPics[1].src

        loadNewImage(firstPics[0])
        loadNewImage(firstPics[1])

        waitingForImages = false
    }

    if (waitingForImages)
        checkIfReady()

    if (loadedImages.length === 5) console.log('loaded at max')
}

function checkIfReady() {
    // we are only waiting for images if literally none are loaded
    // (to replace the back buffer on swap)
    waitingForImages = loadedImages.length === 0
}

function voteCurrentPic(isUpvote) {
    const votedImageSrc = activeBuffer.src
    // TODO: send this (keep the id in a data attr maybe?)

    const [ nextImageElem ] = loadedImages.splice(0, 1)
    activeBuffer = (activeBuffer == mainImg1)
        ? mainImg2
        : mainImg1

    backBuffer = (activeBuffer == mainImg1)
        ? mainImg2
        : mainImg1

    backBuffer.src = nextImageElem.src
    loadNewImage(nextImageElem)
}

function loadNewImage(bufferElem) {
    console.log('loading next')

    fetch('/pictures/next')
        .then(res => res.json())
        .then(res => bufferElem.src = `/static/${res}.jpg`)
        .catch(err => console.log(`error getting next id: ${err}`))
}
