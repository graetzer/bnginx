/**
 * dat.globe Javascript WebGL Globe Toolkit
 * http://dataarts.github.com/dat.globe
 *
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
  var imgDir = opts.imgDir || '/public/globe/';

  var Shaders = {
    /*'earth': {
      uniforms: {
        'texture': {
          type: 't',
          value: null
        }
      },
      vertexShader: [
        'varying vec3 vNormal;',
        'varying vec2 vUv;',
        'void main() {',
        'gl_Position = projectionMatrix * modelViewMatrix * vec4( position, 1.0 );',
        'vNormal = normalize( normalMatrix * normal );',
        'vUv = uv;',
        '}'
      ].join('\n'),
      fragmentShader: [
        'uniform sampler2D texture;',
        'varying vec3 vNormal;',
        'varying vec2 vUv;',
        'void main() {',
        'vec3 diffuse = texture2D( texture, vUv ).xyz;',
        'float intensity = 1.05 - dot( vNormal, vec3( 0.0, 0.0, 1.0 ) );',
        'vec3 atmosphere = vec3( 1.0, 1.0, 1.0 ) * pow( intensity, 3.0 );',
        'gl_FragColor = vec4( diffuse + atmosphere, 1.0 );',
        '}'
      ].join('\n')
    },
    'atmosphere': {
      uniforms: {},
      vertexShader: [
        'varying vec3 vNormal;',
        'void main() {',
        'vNormal = normalize( normalMatrix * normal );',
        'gl_Position = projectionMatrix * modelViewMatrix * vec4( position, 1.0 );',
        '}'
      ].join('\n'),
      fragmentShader: [
        'varying vec3 vNormal;',
        'void main() {',
        'float intensity = pow( 0.8 - dot( vNormal, vec3( 0, 0, 1.0 ) ), 12.0 );',
        'gl_FragColor =  vec4( 0.75 ) * intensity;',
        '}'
      ].join('\n')
    }*/
    /*,
    'locations' : {
      uniforms: {
        'texture': { type: 't', value: null }
      },
      vertexShader: [
        'varying vec2 vUv;',
        'void main() {',
          'gl_Position = projectionMatrix * modelViewMatrix * vec4( position, 1.0 );',
          'vUv = uv;',
        '}'
      ].join('\n'),
      fragmentShader: [
        'uniform sampler2D texture;',
        'varying vec2 vUv;',
        'void main() {',
          'vec3 tex = texture2D( texture, vUv ).xyz;',
          'float a = distance(vUv - vec2(0.5));',
          'float b = 1 - a;',
          'gl_FragColor = vec4(tex * a + vec3(1,0,0)*b, b);',
        '}'
      ].join('\n')
    }*/
  };

  var camera, scene, renderer, w, h;
  var mesh, atmosphere, point;

  var overRenderer, mouseDown;

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
      y: Math.PI / 6.0
    },
    targetOnDown = {
      x: 0,
      y: 0
    };

  var distance = 100000,
    distanceTarget = 100000;
  var padding = 40;
  var PI_HALF = Math.PI / 2;

  function init() {

    var shader, uniforms, material;
    w = opts.width || 600;
    h = opts.height ||600;

    camera = new THREE.PerspectiveCamera(30, w / h, 1, 10000);
    camera.position.z = distance;

    scene = new THREE.Scene();

    var geometry = new THREE.SphereGeometry(200, 40, 30);
    /*shader = Shaders['earth'];
    uniforms = THREE.UniformsUtils.clone(shader.uniforms);
    uniforms['texture'].value = THREE.ImageUtils.loadTexture(imgDir + 'world.jpg');

    material = new THREE.ShaderMaterial({
      uniforms: uniforms,
      vertexShader: shader.vertexShader,
      fragmentShader: shader.fragmentShader
    });*/
    var material = new THREE.MeshBasicMaterial( { color: 0xffffff } );

    mesh = new THREE.Mesh(geometry, material);
    mesh.rotation.y = Math.PI;
    scene.add(mesh);

    // Atmosphere
    /*shader = Shaders['atmosphere'];
    uniforms = THREE.UniformsUtils.clone(shader.uniforms);

    material = new THREE.ShaderMaterial({
      uniforms: uniforms,
      vertexShader: shader.vertexShader,
      fragmentShader: shader.fragmentShader,
      side: THREE.BackSide,
      blending: THREE.AdditiveBlending,
      transparent: true
    });

    mesh = new THREE.Mesh(geometry, material);
    mesh.scale.set(1.1, 1.1, 1.1);
    scene.add(mesh);*/

    // Prototype for the dots
    geometry = new THREE.BoxGeometry(0.75, 0.75, 1);
    geometry.applyMatrix(new THREE.Matrix4().makeTranslation(0, 0, -0.5));
    point = new THREE.Mesh(geometry);

    // Setup renderer
    renderer = new THREE.WebGLRenderer({
      antialias: true
    });
    renderer.setSize(w, h);
    renderer.setClearColor( 0xffffff, 1);
    container.appendChild(renderer.domElement);

    container.addEventListener('mousedown', onMouseDown, false);
    container.addEventListener('mousemove', onMouseMove, false);
    container.addEventListener('mouseup', onMouseUp, false);
    container.addEventListener('mousewheel', onMouseWheel, false);
    container.addEventListener('mouseover', function() {
      overRenderer = true;
    }, false);
    container.addEventListener('mouseout', function() {
      overRenderer = false;
    }, false);

    document.addEventListener('keydown', onDocumentKeyDown, false);

    window.addEventListener('resize', onWindowResize, false);
  }

  function setData(worldMap, locations) {
    var lat, lng, size, i, x, y;
    this.locations = locations;

    var subgeo = new THREE.Geometry();
    /*for (i = 0; i < data.length; i += 2) {
      var color = new THREE.Color();
      color.setRGB(0, 0, 0);
      addPoint(data[i], data[i + 1], 1, color, subgeo);
    }*/
    var color1 = new THREE.Color();
    color1.setRGB(0, 0, 0);

    var color2 = new THREE.Color();
    color2.setRGB(0.5, 0.5, 0.5, 0.8);

    var data = worldMap.data;
    for (i = 0; i < data.length; i += 4 * 42) {// some skip values work better than others
      var intensity = data[i] + data[i+1] + data[i+2];

      x = Math.floor(i / 4) % worldMap.width;
      y = Math.floor(Math.floor(i/4) / worldMap.width);
      lat = 90 - 180 * (y/worldMap.height);// equilateral projection
      lng = 360 * (x/worldMap.width) - 180;
      if (intensity > 0x33 * 3) {
        addPoint(lat, lng, 0.7, 5, color1, subgeo);
      } else {
        addPoint(lat, lng, 0.5, 1, color2, subgeo);
      }
    }

    for (i = 0; i < locations.length; i++) {
      var type = locations[i].type;
      var color = new THREE.Color();
      color.setHSL((0.6 - (type * 0.5)), 1.0, 0.5);
      addPoint(locations[i].lat, locations[i].lng, type, color, subgeo);
    }

    this._baseGeometry = subgeo;
  };

  function createPoints() {
    if (this._baseGeometry !== undefined) {
      this.points = new THREE.Mesh(this._baseGeometry, new THREE.MeshBasicMaterial({
        color: 0xffffff,
        vertexColors: THREE.FaceColors,
        morphTargets: false
      }));

      scene.add(this.points);
    }
  }

  function addPoint(lat, lng, xy, h, color, subgeo) {

    var phi = (90 - lat) * Math.PI / 180;
    var theta = (180 - lng) * Math.PI / 180;

    point.position.x = 200 * Math.sin(phi) * Math.cos(theta);
    point.position.y = 200 * Math.cos(phi);
    point.position.z = 200 * Math.sin(phi) * Math.sin(theta);

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

    mouseOnDown.x = -event.clientX;
    mouseOnDown.y = event.clientY;
    targetOnDown.x = target.x;
    targetOnDown.y = target.y;

    container.style.cursor = 'move';
    mouseDown = true;
  }

  function onMouseMove(event) {
    mouse.x = -event.clientX;
    mouse.y = event.clientY;
    if (mouseDown) {
      var zoomDamp = distance / 1000;
      target.x = targetOnDown.x + (mouse.x - mouseOnDown.x) * 0.005 * zoomDamp;
      target.y = targetOnDown.y + (mouse.y - mouseOnDown.y) * 0.005 * zoomDamp;

      target.y = target.y > PI_HALF ? PI_HALF : target.y;
      target.y = target.y < -PI_HALF ? -PI_HALF : target.y;
    } else {

    }
  }

  function onMouseUp(event) {
    container.style.cursor = 'auto';
    mouseDown = false;
  }

  function onMouseWheel(event) {
    event.preventDefault();
    if (overRenderer) {
      zoom(event.wheelDeltaY * 0.3);
    }
    return false;
  }

  function onDocumentKeyDown(event) {
    switch (event.keyCode) {
      case 38:
        zoom(100);
        event.preventDefault();
        break;
      case 40:
        zoom(-100);
        event.preventDefault();
        break;
    }
  }

  function onWindowResize(event) { // TODO fix window resize
    w = container.offsetWidth || 600;
    h = container.offsetHeight || 600;
    /*camera.aspect = w / h;
    camera.updateProjectionMatrix();
    renderer.setSize( container.offsetWidth, container.offsetHeight );*/
  }

  function zoom(delta) {
    distanceTarget -= delta;
    distanceTarget = distanceTarget > 1000 ? 1000 : distanceTarget;
    distanceTarget = distanceTarget < 350 ? 350 : distanceTarget;
  }

  function animate() {
    requestAnimationFrame(animate);
    render();
  }

  function render() {
    zoom(curZoomSpeed);// modifies distanceTarget

    rotation.x += (target.x - rotation.x) * 0.1;
    rotation.y += (target.y - rotation.y) * 0.1;
    distance += (distanceTarget - distance) * 0.3;

    camera.position.x = distance * Math.sin(rotation.x) * Math.cos(rotation.y);
    camera.position.y = distance * Math.sin(rotation.y);
    camera.position.z = distance * Math.cos(rotation.x) * Math.cos(rotation.y);

    camera.lookAt(mesh.position);

    renderer.render(scene, camera);
  }

  init();
  this.animate = animate;

  this.setData = setData;
  this.createPoints = createPoints;
  this.renderer = renderer;
  this.scene = scene;

  return this;
};
