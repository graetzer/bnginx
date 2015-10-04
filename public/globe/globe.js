/**
 * Globe Animation
 *
 * Copyright 2015 Simon P. Grätzer
 * Copyright 2011 Data Arts Team, Google Creative Lab
 *
 * Licensed under the Apache License, Version 2.0 (the 'License');
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 */

var DAT = DAT || {};

DAT.Globe = function(container, opts) {
  opts = opts || {};

  var camera, scene, renderer, w, h;
  var mesh, atmosphere, point;

  var overRenderer, mouseDown, zoomTouchSpread;

  var curZoomSpeed = 0;
  var zoomSpeed = 50;

  var mouse = {
      x: 0,
      y: 0
    },
    mouseOnDown = {
      x: 0,
      y: 0
    };
  var rotation = {
      x: 0,
      y: 0
    },
    target = {
      x: Math.PI * 3 / 2,
      y: Math.PI / 8.0
    },
    targetOnDown = {
      x: 0,
      y: 0
    };

  var distance = 100000,
    distanceTarget = 100000;
  var PI_HALF = Math.PI / 2;

  function init() {

    var shader, uniforms, material;
    //600 is enough to get the globe in full and have a good render performance
    w = container.clientWidth || 600;
    h = Math.min(600, w) || 600;

    camera = new THREE.PerspectiveCamera(25, w / h, 1, 10000);
    camera.position.z = distance;

    scene = new THREE.Scene();

    var geometry = new THREE.SphereGeometry(200, 40, 30);
    var material = new THREE.MeshBasicMaterial({
      color: 0xffffff
    });

    mesh = new THREE.Mesh(geometry, material);
    mesh.rotation.y = Math.PI;
    scene.add(mesh);

    // Prototype for the dots
    geometry = new THREE.BoxGeometry(0.75, 0.75, 1);
    geometry.applyMatrix(new THREE.Matrix4().makeTranslation(0, 0, -0.5));
    point = new THREE.Mesh(geometry);

    // Setup renderer
    renderer = new THREE.WebGLRenderer({
      antialias: true
    });
    renderer.setSize(w, h);
    renderer.setClearColor(0xffffff, 1);
    renderer.domElement.style.width = "100%";
    renderer.domElement.style.height = "auto";
    container.appendChild(renderer.domElement);

    container.addEventListener('mousedown', onMouseDown, false);
    container.addEventListener('mousemove', onMouseMove, false);
    container.addEventListener('mouseup', onMouseUp, false);
    container.addEventListener('wheel', onMouseWheel, false);
    container.addEventListener("touchstart", onMouseDown, false);
    container.addEventListener("touchend", onMouseUp, false);
    container.addEventListener("touchmove", onMouseMove, false);
    container.addEventListener('mouseover', function() {
      overRenderer = true;
    }, false);
    container.addEventListener('mouseout', function() {
      overRenderer = false;
      mouseDown = false;
    }, false);

    document.addEventListener('keydown', onDocumentKeyDown, false);
    //window.addEventListener('resize', onWindowResize, false);
  }

  function setData(worldMap, places) {
    var lat, lng, step, i, x, y, intensity;

    var subgeo = new THREE.Geometry();
    var color1 = new THREE.Color(0, 0, 0);
    var color2 = new THREE.Color(0.5, 0.5, 0.5, 0.8);

    var data = worldMap.data;

    var xx = 0;
    for (y = 0; y < worldMap.height; y += 4) {
      //step = Math.floor(35 * Math.pow(2 * y / worldMap.height - 1, 4) + 4);
      ys = y / worldMap.height;
      x = 0;
      if (0.128 < ys && ys < 0.76) step = 5;
      else {
        if (ys <= 0.128) {// 2^(3 + distBorder) at least 8, doubles every 5 px
          step = Math.floor(Math.pow(2, (0.128 * worldMap.height - y) / 7) + 3);
        } else {
          step = Math.floor(Math.pow(2, (y - 0.76*worldMap.height) / 7) + 3);
        }
        x = Math.floor(step);
      }
      for (; x < worldMap.width; x += step) {
        i = (y * worldMap.width + x) * 4;
        intensity = data[i] + data[i + 1] + data[i + 2];

        lat = 90 - 180 * (y / worldMap.height); // equilateral projection
        lng = 360 * (x / worldMap.width) - 180;
        if (intensity > 0x33 * 3) {
          addPoint(lat, lng, 0.8, 5, color1, subgeo);
        } else {
          addPoint(lat, lng, 0.5, 1, color2, subgeo);
        }
        xx++;
      }
    }
    addPoint(90, 0, 0.5, 1, color2, subgeo);// northpole
    addPoint(-90, 0, 0.5, 1, color2, subgeo);//southpole
    console.log("%d nodes", xx+2);
    this.points = new THREE.Mesh(subgeo, new THREE.MeshBasicMaterial({
      color: 0xffffff,
      vertexColors: THREE.FaceColors,
      morphTargets: false
    }));
    scene.add(this.points);

    // Add the Places
    for (i = 0; i < places.length; i++) {
      var vertex = new THREE.Vector3();
      convertLatLng(places[i].lat, places[i].lng, vertex, 205);

      var sprite = THREE.ImageUtils.loadTexture(places[i].coverUrl);
      sprite.minFilter = THREE.LinearFilter;
      var material = new THREE.PointsMaterial({
        size: 15,
        map: sprite,
        alphaTest: 0.5,
        transparent: true
      });
      var particleGeo = new THREE.Geometry();
      particleGeo.vertices.push(vertex);
      var particles = new THREE.Points(particleGeo, material);
      scene.add(particles);
    }
  };

  function convertLatLng(lat, lng, outXYZ, scale) {
    var phi = (90 - lat) * Math.PI / 180;
    var theta = (180 - lng) * Math.PI / 180;

    outXYZ.x = scale * Math.sin(phi) * Math.cos(theta);
    outXYZ.y = scale * Math.cos(phi);
    outXYZ.z = scale * Math.sin(phi) * Math.sin(theta);
  }

  function addPoint(lat, lng, xy, h, color, subgeo) {
    convertLatLng(lat, lng, point.position, 200);
    point.lookAt(mesh.position);

    point.scale.x = xy;
    point.scale.y = xy;
    point.scale.z = h; // avoid non-invertible matrix
    point.updateMatrix();

    for (var i = 0; i < point.geometry.faces.length; i++) {
      point.geometry.faces[i].color = color;
    }
    if (point.matrixAutoUpdate) {
      point.updateMatrix();
    }
    subgeo.merge(point.geometry, point.matrix);
  }

  function onMouseDown(event) {
    event.preventDefault();
    if (event.targetTouches) {
      var tts = event.targetTouches;
      if (tts.length == 2) {// pinch to zoom
        var xx = tts[0].clientX - tts[1].clientX;
        var yy = tts[0].clientY - tts[1].clientY;
        zoomTouchSpread = Math.sqrt(xx*xx  + yy*yy);
        return;
      }
      event = event.targetTouches[0];// has clientXY
    }

    mouseOnDown.x = -event.clientX;
    mouseOnDown.y = event.clientY;
    targetOnDown.x = target.x;
    targetOnDown.y = target.y;

    container.style.cursor = 'move';
    mouseDown = true;
  }

  function onMouseMove(event) {
    if (event.targetTouches) {
      var tts = event.targetTouches;
      if (tts.length == 2) {// pinch to zoom
        var xx = tts[0].clientX - tts[1].clientX;
        var yy = tts[0].clientY - tts[1].clientY;
        var zoomDelta = Math.sqrt(xx*xx  + yy*yy) - zoomTouchSpread;
        zoom(zoomDelta * 0.4);// 0.4 feels smooth on an iphone 6
        return;
      }
      event = event.targetTouches[0];// has clientXY
    }

    mouse.x = -event.clientX;
    mouse.y = event.clientY;
    if (mouseDown) {
      var zoomDamp = distance / 1000;
      target.x = targetOnDown.x + (mouse.x - mouseOnDown.x) * 0.005 * zoomDamp;
      target.y = targetOnDown.y + (mouse.y - mouseOnDown.y) * 0.005 * zoomDamp;

      target.y = target.y > PI_HALF ? PI_HALF : target.y;
      target.y = target.y < -PI_HALF ? -PI_HALF : target.y;
    }
  }

  function onMouseUp(event) {
    container.style.cursor = 'auto';
    mouseDown = false;
  }

  function onMouseWheel(event) {
    event.preventDefault();
    if (overRenderer) {
      zoom(-event.deltaY * 0.3);
    }
    return false;
  }

  function onDocumentKeyDown(event) {
    var key = event.which || event.keyCode || 0;
    switch (key) {
      case 38: // up arrow
        zoom(100);
        event.preventDefault();
        break;
      case 40: // down arrow
        zoom(-100);
        event.preventDefault();
        break;
      case 87: // w
        target.y += 0.1;
        event.preventDefault();
        break;
      case 83: // s
        target.y -= 0.1;
        event.preventDefault();
        break;
      case 65: // a
      case 37: //left arrow
        target.x -= 0.1;
        event.preventDefault();
        break;
      case 68: // d
      case 39: //right arrow
        target.x += 0.1;
        event.preventDefault();
        break;
    }
    target.y = target.y > PI_HALF ? PI_HALF : target.y;
    target.y = target.y < -PI_HALF ? -PI_HALF : target.y;
  }

  function onWindowResize(event) { // TODO fix window resize
    w = container.clientWidth || 600;
    h = Math.min(600, w) || 600;
    camera.aspect = w / h;
    camera.updateProjectionMatrix();
    renderer.setSize( w, h );
  }

  function zoom(delta) {
    distanceTarget -= delta;
    distanceTarget = distanceTarget > 1000 ? 1000 : distanceTarget;
    distanceTarget = distanceTarget < 275 ? 275 : distanceTarget;
  }

  function animate() {
    render();
    requestAnimationFrame(animate);
  }

  function render() {
    zoom(curZoomSpeed); // modifies distanceTarget

    rotation.x += (target.x - rotation.x) * 0.1;
    rotation.y += (target.y - rotation.y) * 0.1;
    distance += (distanceTarget - distance) * 0.2;

    camera.position.x = distance * Math.sin(rotation.x) * Math.cos(rotation.y);
    camera.position.y = distance * Math.sin(rotation.y);
    camera.position.z = distance * Math.cos(rotation.x) * Math.cos(rotation.y);
    camera.lookAt(mesh.position);

    renderer.render(scene, camera);
  }

  function showLatLng(lat, lng) { // No idea why this works
    target.x = lng * Math.PI / 180 - PI_HALF;
    target.y = lat * Math.PI / 180;
    distanceTarget = 500;
  }

  init();
  //support different inital position
  if (opts.lat && opts.lng) {
    showLatLng(opts.lat, opts.lng);
    distanceTarget = 50000; // don't zoom in quite so far
  }
  this.animate = animate;

  this.setData = setData;
  this.renderer = renderer;
  this.scene = scene;
  this.showLatLng = showLatLng;

  return this;
};
