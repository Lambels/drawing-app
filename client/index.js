TYPEFIRST = 3;
TYPECLEAR = 4;
TYPEPOINT = 7;
TYPEPOINTSTART = 8;
TYPEPOINTEND = 9;

window.onload = () => {
  const draw = (x, y, color) => {
    if (!isDrawing) {
      return;
    }

    ctx.lineWidth = 4;
    ctx.lineCap = "round";
    ctx.strokeStyle = color;

    ctx.lineTo(x, y);
    ctx.stroke();
  };

  const canvas = document.getElementById("drawing-board");
  canvas.width = window.innerWidth;
  canvas.height = window.innerHeight;
  var ctx = canvas.getContext("2d");

  canvas.addEventListener("mousemove", (e) => {
    draw(e.clientX, e.clientY, strokeColor);
    // if drawing state send line info to other clients.
    if (isDrawing) {
      ws.send(
        JSON.stringify({
          type: TYPEPOINT,
          id: id,
          color: strokeColor,
          x: e.clientX,
          y: e.clientY,
        })
      );
    }
  });

  canvas.addEventListener("mousedown", (e) => {
    isDrawing = true;
    startX = e.clientX;
    startY = e.clientY;
    // send starting point to other clients.
    ws.send(
      JSON.stringify({
        type: TYPEPOINTSTART,
        id: id,
        color: strokeColor,
        x: startX,
        y: startY,
      })
    );
  });

  canvas.addEventListener("mouseup", (e) => {
    isDrawing = false;
    ctx.stroke();
    ctx.beginPath();
    ws.send(
      JSON.stringify({
        type: TYPEPOINTEND,
        id: id,
        color: strokeColor,
        x: startX,
        y: startY,
      })
    );
  });

  var id = -1;
  var strokeColor = "";
  var isDrawing = false;
  let startX;
  let startY;

  const ws = new WebSocket("ws://localhost:8080/ws");
  ws.onmessage = (e) => {
    payload = JSON.parse(e.data);
    switch (payload.type) {
      case TYPEFIRST:
        // sync server side and client side with appropiate data.
        id = payload.id;
        strokeColor = payload.color;
        break;

      case TYPECLEAR:
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        break;

      case TYPEPOINT:
        draw(payload.x, payload.y, payload.color);
        break;

      case TYPEPOINTSTART:
        isDrawing = true;
        break;

      case TYPEPOINTEND:
        isDrawing = false;
        ctx.stroke();
        ctx.beginPath();
        break;
    }
  };

  ws.onclose = (e) => {
    console.log("closed!", e);
  };
};
