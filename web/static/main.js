function onload() {
    var type = "WebGL"
    if (!PIXI.utils.isWebGLSupported()) {
        type = "canvas"
    }

    PIXI.utils.sayHello(type)

    var app = new PIXI.Application({
        width: 512,
        height: 512,
        resolution: window.devicePixelRatio,
        autoDensity: true,
    })
    
    document.body.appendChild(app.view)

    let points = [0, 42, 1, 42, 1, 43]
    
    let rect = new PIXI.Graphics()
    rect.lineStyle(2, 0xff3300, 1)
    //rect.beginFill(0x66ccff)
    rect.drawPolygon(points)
    //rect.endFill()
    rect.x = 0
    rect.y = 0

    app.stage.addChild(rect)

    var sock = new WebSocket("ws://localhost:8000/ws/")
}
