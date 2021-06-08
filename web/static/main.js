var pixi, sock, polygons

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
        resolution: window.devicePixelRatio,
        autoDensity: true,
    })
    
    document.body.appendChild(pixi.view)

    /* let points = [0, 42, 1, 42, 1, 43]
    
    let rect = new PIXI.Graphics()
    rect.lineStyle(2, 0xff3300, 1)
    //rect.beginFill(0x66ccff)
    rect.drawPolygon(points)
    //rect.endFill()
    rect.x = 0
    rect.y = 0

    app.stage.addChild(rect)
    */
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
        poly.lineStyle(2, 0xff3300, 1)
        poly.drawPolygon(polygons[i])
        poly.x = 0
        poly.y = 0
        pixi.stage.addChild(poly)
    }
}
