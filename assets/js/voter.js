const mainImg1 = document.getElementById('main-1')
const mainImg2 = document.getElementById('main-2')

let activeBuffer = mainImg1

let loadedImages = []
const bufferElements = document.getElementsByClassName('preload')

const waitingForImages = true

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
        addLoadedImage({
            ref: bufferElements[i],
            src: bufferElements[i].src,
        })
    })
}

function addLoadedImage(loaded) {
    loadedImages.push(loaded)

    if (waitingForImages)
        checkIfReady()
}

function checkIfReady() {
    // uhh this is incorrect
    waitingForImages = loadedImages.length < 5
}

function voteCurrentPic(isUpvote) {
    const votedImageSrc = activeBuffer.src

    // find img element that is now obsolete
    const oldBufferIndex = loadedImages.findIndex(
        it => it.src === votedImageSrc,
    )

    const removed = loadedImages.splice(oldBufferIndex, 1)
    if (removed.length !== 1) console.error('this should not happen')

    loadNewImage(oldBuffer)

    // set new active buffer
    activeBuffer = (activeBuffer == mainImg1)
        ? mainImg2
        : mainImg1
}

function loadNewImage(bufferElem) {
    console.log(`resetting image ${bufferElem.src}`)
}
