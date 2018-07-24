import Hammer from 'hammerjs'

const _frontBuffer = document.getElementById('main-1')
const _backBuffer = document.getElementById('main-2')

const _frontBufferMc = new Hammer(_frontBuffer)
const _backBufferMc = new Hammer(_backBuffer)

const loadingOverlay = document.getElementById('loading-overlay')

const buffers = [
    {
        ref: _frontBuffer,
        mc: _frontBufferMc,
        id: '',
    },
    {
        ref: _backBuffer,
        mc: _backBufferMc,
        id: '',
    },
]

let frontBuffer = 0

let loadedImages = []
let initialized = false
const bufferElements = document.getElementsByClassName('preload')

let waitingForImages = true

window.addEventListener('keyup', ev => {
    if (waitingForImages) return

    if (ev.keyCode === 65) {
        voteCurrentPic(true)
    }

    if (ev.keyCode === 83) {
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

        const [ front, back ] = loadedImages.splice(0, 2)
        const frontId = front.getAttribute('data-pic-id')
        const backId = back.getAttribute('data-pic-id')

        getFrontBuffer().ref.classList.add('front')
        getBackBuffer().ref.classList.add('back')

        setFrontBuffer(front, frontId)
        setBackBuffer(back, backId)

        ;[front, back].forEach(it => loadNewImage(it))

        setWaiting(false)
    }

    if (waitingForImages)
        checkIfReady()
}

function checkIfReady() {
    // we are only waiting for images if literally none are loaded
    // (to replace the back buffer on swap)
    setWaiting(loadedImages.length === 0)
}

function voteCurrentPic(isUpvote) {
    const votedImageSrc = getFrontBuffer().ref.src
    // TODO: send this (keep the id in a data attr maybe?)

    const id = getFrontBuffer().id
    fetch(`/pictures/${id}/votes`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ isUpvote }),
    })
        .then(res => {})
        .catch(err => console.error(`error voting: ${err}`))

    const [ nextImageElem ] = loadedImages.splice(0, 1)
    const nextImageId = nextImageElem.getAttribute('data-pic-id')

    swapBuffers()
    setBackBuffer(nextImageElem, nextImageId)
    loadNewImage(nextImageElem)

    checkIfReady()
    if (waitingForImages) {} // TODO: Check if we even ever need this
}

function loadNewImage(bufferElem) {
    fetch('/pictures/next')
        .then(res => res.json())
        .then(res => {
            bufferElem.src = `/static/${res}.jpg`
            bufferElem.setAttribute('data-pic-id', res)
            checkIfReady()
        })
        .catch(err => console.log(`error getting next id: ${err}`))
}

function getFrontBuffer() { return buffers[frontBuffer] }
function getBackBuffer() { return buffers[1 - frontBuffer] }

function setFrontBuffer(elem, id) {
    buffers[frontBuffer].ref.src = elem.src
    buffers[frontBuffer].id = id
}

function setBackBuffer(elem, id) {
    buffers[1 - frontBuffer].ref.src = elem.src
    buffers[1 - frontBuffer].id = id
}

function swapBuffers() {
    getFrontBuffer().ref.classList.remove('front')
    getFrontBuffer().ref.classList.add('back')

    getBackBuffer().ref.classList.remove('back')
    getBackBuffer().ref.classList.add('front')

    frontBuffer = 1 - frontBuffer
}

function setWaiting(isWaiting) {
    waitingForImages = isWaiting

    if (isWaiting) loadingOverlay.classList.add('active')
    else loadingOverlay.classList.remove('active')
}
