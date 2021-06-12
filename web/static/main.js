var pixi, sock, polygons, centroids

function onload() {
    init()
}

function init() {
    sock = new WebSocket("ws://localhost:8000/ws/")
    sock.onopen = wsOpened
    sock.onmessage = wsGotMessage
    sock.onclose = wsClosed
    sock.onerror = wsError
    
    var type = "WebGL"
    if (!PIXI.utils.isWebGLSupported()) {
        type = "canvas"
    }

    PIXI.utils.sayHello(type)

    pixi = new PIXI.Application({
        width: 600,
        height: 600,
        backgroundColor: 0xff0000,
        resolution: window.devicePixelRatio,
        autoDensity: true,
    })
    
    document.body.appendChild(pixi.view)
}

function wsOpened(evt) {
    sock.send("GEN 512 512 20 10")
}

function wsGotMessage(evt) {
    msg = evt.data
    
    if (msg.startsWith("GENERATED")) {
        sock.send("POL")
    } else if (msg.startsWith("POLYGONS")) {
        polygons = JSON.parse(msg.substring(8))
        console.log("got polygons")
        sock.send("CEN")
    } else if (msg.startsWith("CENTROIDS")) {
        centroids = JSON.parse(msg.substring(9))
        console.log("got centroids")
        render()
    } else if (msg.startsWith("NOGAME")) {
        console.error("no game in progress")
    }
}

function wsClosed(evt) {
    
}

function wsError(evt) {
    console.error(evt)
}

function render() {
    for (var i = 0; i < polygons.length; i++) {
        let poly = new PIXI.Graphics()
        poly.beginFill(0xffffff)
        poly.lineStyle(3, 0x000000, 1)
        poly.drawPolygon(polygons[i])
        poly.endFill()
        
        poly.x = 0
        poly.y = 0
        pixi.stage.addChild(poly)
    }

    for (var i = 0; i < centroids.length; i++) {
        let center = new PIXI.Graphics()
        center.beginFill(0x000000)
        center.lineStyle(2, 0xffffff, 1)
        center.drawRoundedRect(-12.5, -15, 25, 30, 4)
        center.endFill()
        
        center.x = centroids[i].x
        center.y = centroids[i].y
        pixi.stage.addChild(center)
    }
}
