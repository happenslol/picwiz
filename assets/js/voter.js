import Hammer from 'hammerjs'

const _frontBuffer = document.getElementById('main-1')
const _backBuffer = document.getElementById('main-2')

const _frontBufferMc = new Hammer(_frontBuffer)
const _backBufferMc = new Hammer(_backBuffer)

const voteOverlay = document.getElementById('vote-overlay')
const upvoteOverlay = document.getElementById('up-overlay')
const downvoteOverlay = document.getElementById('down-overlay')

const loadingOverlay = document.getElementById('loading-overlay')

let upvoteTriggered = false
let downvoteTriggered = false

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

_frontBufferMc.add(new Hammer.Pan({
    direction: Hammer.DIRECTION_HORIZONTAL,
    threshold: 0,
}))

_backBufferMc.add(new Hammer.Pan({
    direction: Hammer.DIRECTION_HORIZONTAL,
    threshold: 0,
}))

_frontBuffer.addEventListener('dragstart', ev => ev.preventDefault())
_backBuffer.addEventListener('dragstart', ev => ev.preventDefault())

_frontBufferMc.on('pan', ev => handleBufferPan(0)(ev))
_backBufferMc.on('pan', ev => handleBufferPan(1)(ev))

_frontBufferMc.on('panend', ev => handleBufferPanEnd(0)(ev))
_backBufferMc.on('panend', ev => handleBufferPanEnd(1)(ev))

let frontBuffer = 0

function handleBufferPan(bufferId) {
    if (frontBuffer !== bufferId) return () => {}
    return ev => {
        const screenWidth = document.documentElement.clientWidth
        let percentMoved = (ev.deltaX / screenWidth) * 2
        let absMoved = Math.abs(percentMoved)
        absMoved = absMoved > 1.0 ? 1.0 : absMoved
        const t = 0.7 + (0.25 * absMoved)
        const fb = getFrontBuffer()
        const bb = getBackBuffer()

        voteOverlay.classList.add('active')

        fb.ref.style.transform = `translate3d(${ev.deltaX}px, 0, 0)`
        bb.ref.style.transform = `scale3d(${t}, ${t}, 1.0)`

        if (percentMoved > 0.35) {
            downvoteOverlay.classList.add('active')
        } else if (percentMoved < -0.35) {
            upvoteOverlay.classList.add('active')
        } else {
            upvoteOverlay.classList.remove('active')
            downvoteOverlay.classList.remove('active')
        }
    }
}

function handleBufferPanEnd(bufferId) {
    if (frontBuffer !== bufferId) return () => {}
    return ev => {
        const fb = getFrontBuffer()
        const bb = getBackBuffer()

        fb.ref.classList.add('animating')
        bb.ref.classList.add('animating')

        console.dir(ev)

        fb.ref.style = ''
        bb.ref.style = ''

        voteOverlay.classList.remove('active')

        requestAnimationFrame(() => {
            upvoteOverlay.classList.remove('active')
            downvoteOverlay.classList.remove('active')
        })

        setTimeout(() => {
            fb.ref.classList.remove('animating')
            bb.ref.classList.remove('animating')
        }, 260)
    }
}

let loadedImages = []
let initialized = false
const bufferElements = document.getElementsByClassName('preload')

let waitingForImages = true
let animating = false

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
