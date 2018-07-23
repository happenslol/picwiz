import Hammer from 'hammerjs'

const _frontBuffer = document.getElementById('main-1')
const _backBuffer = document.getElementById('main-2')

const _frontBufferMc = new Hammer(_frontBuffer)
const _backBufferMc = new Hammer(_backBuffer)

const buffers = [
    {
        ref: _frontBuffer,
        mc: _frontBufferMc,
    },
    {
        ref: _backBuffer,
        mc: _backBufferMc,
    },
]

let frontBuffer = 0

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

        initialized = true

        const [ active, back ] = loadedImages.splice(0, 2)

        getFrontBuffer().ref.classList.add('front')
        getBackBuffer().ref.classList.add('back')

        setFrontBuffer(active)
        setBackBuffer(back)

        ;[active, back].forEach(it => loadNewImage(it))

        waitingForImages = false
    }

    if (waitingForImages)
        checkIfReady()
}

function checkIfReady() {
    // we are only waiting for images if literally none are loaded
    // (to replace the back buffer on swap)
    waitingForImages = loadedImages.length === 0
}

function voteCurrentPic(isUpvote) {
    const votedImageSrc = getFrontBuffer().ref.src
    // TODO: send this (keep the id in a data attr maybe?)

    const [ nextImageElem ] = loadedImages.splice(0, 1)
    swapBuffers()
    setBackBuffer(nextImageElem)
    loadNewImage(nextImageElem)

    checkIfReady()
    if (waitingForImages) {} // TODO: Check if we even ever need this
}

function loadNewImage(bufferElem) {
    fetch('/pictures/next')
        .then(res => res.json())
        .then(res => {
            bufferElem.src = `/static/${res}.jpg`
            checkIfReady()
        })
        .catch(err => console.log(`error getting next id: ${err}`))
}

function getFrontBuffer() { return buffers[frontBuffer] }
function getBackBuffer() { return buffers[1 - frontBuffer] }
function setFrontBuffer(elem) { buffers[frontBuffer].ref.src = elem.src }
function setBackBuffer(elem) { buffers[1 - frontBuffer].ref.src = elem.src }
function swapBuffers() {
    getFrontBuffer().ref.classList.remove('front')
    getFrontBuffer().ref.classList.add('back')

    getBackBuffer().ref.classList.remove('back')
    getBackBuffer().ref.classList.add('front')

    frontBuffer = 1 - frontBuffer
}
