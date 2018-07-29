import Hammer from 'hammerjs'

const _frontBuffer = document.getElementById('main-1')
const _backBuffer = document.getElementById('main-2')

const voteOverlay = document.getElementById('vote-overlay')
const upvoteOverlay = document.getElementById('up-overlay')
const downvoteOverlay = document.getElementById('down-overlay')

const loadingOverlay = document.getElementById('loading-overlay')

const buffers = [
    {
        ref: _frontBuffer,
        id: '',
    },
    {
        ref: _backBuffer,
        id: '',
    },
]

_frontBuffer.addEventListener('dragstart', ev => ev.preventDefault())
_backBuffer.addEventListener('dragstart', ev => ev.preventDefault())

const mc = new Hammer(document.body)

mc.add(new Hammer.Pan({
    direction: Hammer.DIRECTION_HORIZONTAL,
    threshold: 0,
}))

// mc.add(new Hammer.Pinch({
    // threshold: 0,
// }))

mc.on('panleft panright', handleBufferPan)
mc.on('panend pancancel', handleBufferPanEnd)
// mc.on('pinch', handleBufferPinch)

document.body.addEventListener('click touch', ev => {
    ev.preventDefault()
    ev.stopPropagation()
})

let frontBuffer = 0

function handleBufferPan(ev) {
    const screenWidth = document.documentElement.clientWidth
    let percentMoved = (ev.deltaX / screenWidth) * 2
    let absMoved = Math.abs(percentMoved)
    absMoved = absMoved > 1.0 ? 1.0 : (Math.floor(absMoved * 100) / 100)
    const t = 0.7 + (0.25 * absMoved)
    const o = 0.9 * absMoved
    const fb = getFrontBuffer()
    const bb = getBackBuffer()

    voteOverlay.classList.add('active')

    fb.ref.style.transform = `translate3d(${ev.deltaX}px, 0, 0)`
    bb.ref.style.transform = `scale3d(${t}, ${t}, 1.0)`
    bb.ref.style.opacity = `${o}`

    if (percentMoved > 0.5) {
        downvoteOverlay.classList.add('active')
    } else if (percentMoved < -0.5) {
        upvoteOverlay.classList.add('active')
    } else {
        upvoteOverlay.classList.remove('active')
        downvoteOverlay.classList.remove('active')
    }
}

function handleBufferPanEnd(ev) {
    const fb = getFrontBuffer()
    const bb = getBackBuffer()

    const screenWidth = document.documentElement.clientWidth
    let percentMoved = (ev.deltaX / screenWidth) * 2
    let absMoved = Math.abs(percentMoved)
    absMoved = absMoved > 1.0 ? 1.0 : (Math.floor(absMoved * 100) / 100)
    const t = 0.7 + (0.25 * absMoved)

    if (percentMoved > 0.5) {
        voteCurrentPic(false)
    } else if (percentMoved < -0.5) {
        voteCurrentPic(true)
    } else {
        fb.ref.classList.add('animating')
        bb.ref.classList.add('animating')
        document.body.classList.add('animating')

        fb.ref.style = ''
        bb.ref.style = ''

        setTimeout(() => {
            requestAnimationFrame(() => {
                fb.ref.classList.remove('animating')
                bb.ref.classList.remove('animating')
                document.body.classList.remove('animating')
            })
        }, 150)
    }

    voteOverlay.classList.remove('active')
    requestAnimationFrame(() => {
        upvoteOverlay.classList.remove('active')
        downvoteOverlay.classList.remove('active')
    })
}

function handleBufferPinch(ev) {
    ev.preventDefault()

    console.dir(ev)
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
    if (bufferElements[i].complete)
        addLoadedImage(bufferElements[i])

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
    const fb = getFrontBuffer()
    const bb = getBackBuffer()
    const votedImageSrc = fb.ref.src

    fb.ref.classList.add('animating')
    bb.ref.classList.add('animating')
    document.body.classList.add('animating')

    bb.ref.style.transform = 'scale3d(1.0, 1.0, 1.0)'
    bb.ref.style.opacity = '1'

    if (isUpvote) {
        fb.ref.style.transform = 'translate3d(-100%, 0, 0)'
    } else {
        fb.ref.style.transform = 'translate3d(100%, 0, 0)'
    }

    setTimeout(() => {
        requestAnimationFrame(() => {
            fb.ref.classList.remove('animating')
            bb.ref.classList.remove('animating')
            document.body.classList.remove('animating')

            getFrontBuffer().ref.classList.remove('front')
            getBackBuffer().ref.classList.remove('back')

            fb.ref.style = ''
            bb.ref.style = ''

            const [ nextImageElem ] = loadedImages.splice(0, 1)
            const nextImageId = nextImageElem.getAttribute('data-pic-id')

            frontBuffer = 1 - frontBuffer

            getFrontBuffer().ref.classList.add('front')
            getBackBuffer().ref.classList.add('back')

            requestAnimationFrame(() => {
                setBackBuffer(nextImageElem, nextImageId)
                loadNewImage(nextImageElem)
            })

            checkIfReady()
        })
    }, 150)

    const id = fb.id
    fetch(`/pictures/${id}/votes`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ isUpvote }),
    })
        .then(res => {})
        .catch(err => console.error(`error voting: ${err}`))

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

function setWaiting(isWaiting) {
    waitingForImages = isWaiting

    if (isWaiting) loadingOverlay.classList.add('active')
    else loadingOverlay.classList.remove('active')
}
