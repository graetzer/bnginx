/**
 * Globe Animation
 *
 * Copyright 2020 Simon P. Grätzer
 * Copyright 2011 Data Arts Team, Google Creative Lab
 *
 * Licensed under the Apache License, Version 2.0 (the 'License');
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 */
"use strict";

var DAT = DAT || {};

DAT.Globe = function(container, worldMap, opts) {
    opts = opts || {};
    opts.places = opts.places || [];

    var camera, scene, renderer;
    var globeMesh, point;
    var mouseOver, mouseDown, zoomTouchSpread;
    var originPos = new THREE.Vector3(0, 0, 0);
    var mouse = { x: 0, y: 0 }, mouseOnDown = { x: 0, y: 0 };
    var rotation = { x: 0, y: 0 }, target = { x: Math.PI * 3 / 2, y: Math.PI / 8.0 }, targetOnDown = { x: 0, y: 0 };
    var distance = 100000, distanceTarget = 100000;
    var PI_HALF = Math.PI / 2;

    function init() {

        var shader, uniforms, material;
        //600 is enough to get the globe in full and have a good render performance
        var w = opts.width || container.clientWidth || 600;
        var h = opts.height || Math.min(600, w) || 600;

        camera = new THREE.PerspectiveCamera(opts.angle || 25, w / h, 1, 10000);
        camera.position.z = distance;
        scene = new THREE.Scene();

        // Setup renderer
        renderer = new THREE.WebGLRenderer({ antialias: true });
        renderer.setSize(w, h);
        renderer.setClearColor(0xffffff, 1);
        renderer.domElement.style.width = opts.width ? opts.width + "px" : "100%";
        renderer.domElement.style.height = opts.height ? opts.height + "px" : "auto";
        container.appendChild(renderer.domElement);

        // Add all event listeners
        if (!opts.disableMouse) {
            container.addEventListener('mousedown', onMouseDown, false);
            container.addEventListener('mousemove', onMouseMove, false);
            container.addEventListener('mouseup', onMouseUp, false);
        }
        if (!opts.disableMouseZoom) container.addEventListener('wheel', onMouseWheel, false);
        if (!opts.disableTouch) {
            container.addEventListener("touchstart", onMouseDown, false);
            container.addEventListener("touchend", onMouseUp, false);
            container.addEventListener("touchmove", onMouseMove, false);
        }
        container.addEventListener('mouseover', function() {
            mouseOver = true;
        }, false);
        container.addEventListener('mouseout', function() {
            mouseOver = false;
            mouseDown = false;
        }, false);
        document.addEventListener('keydown', onDocumentKeyDown, false);
        //window.addEventListener('resize', onWindowResize, false);
        
        // Prototype for the dots
        var geometry = new THREE.BoxGeometry(0.75, 0.75, 1);
        geometry.applyMatrix(new THREE.Matrix4().makeTranslation(0, 0, -0.5));
        point = new THREE.Mesh(geometry);

        //var globeTexData = new Uint8Array(worldMap.width*worldMap.height*3);
        //globeTexData.fill(0xFF);
        var data = worldMap.data;
        var subgeo = new THREE.Geometry();
        var color1 = new THREE.Color(opts.color || "#000");
        var color2 = new THREE.Color(0x666666);
        var step, x, y, nodeCount = 0;
        for (y = 0; y < worldMap.height; y += 3) {
            //step = Math.floor(35 * Math.pow(2 * y / worldMap.height - 1, 4) + 4);
            var yy = y / worldMap.height;// now use as scaled y coordinate

            if (yy <= 0.07) continue; // don't care about arctica, it's not on my map
            else if (yy <= 0.14) step = Math.floor(Math.pow(2, (0.14 * worldMap.height - y) / 6) + 3);
            else if (yy < 0.83) step = 3;
            else break;// Don't care about antarctica, it's not on my map

            for (x = 0; x < worldMap.width; x += step) {
                var i = (y * worldMap.width + x) * 4;
                var intensity = data[i] + data[i + 1] + data[i + 2];
                var lat = 90 - 180 * (y / worldMap.height); // equilateral projection
                var lng = 360 * (x / worldMap.width) - 180;
                if (intensity > 0x33 * 3) {
                    addPoint(lat, lng, 0.8, 7, color1, subgeo);
                    nodeCount++;
                } else if (x % 4 == 0 && y % 4 == 0) {
                    addPoint(lat, lng, 0.8, 1, color2, subgeo); nodeCount++;
                    //i = (y * worldMap.width + x) * 3;
                    //globeTexData.fill(0x66, i, i + 3);
                }
            }
        }
        addPoint(90, 0, 0.5, 1, color2, subgeo);// northpole
        addPoint(-90, 0, 0.5, 1, color2, subgeo);//southpole
        console.log("%d nodes", nodeCount + 2);

        var points = new THREE.Mesh(subgeo, new THREE.MeshBasicMaterial({
            color: 0xffffff,
            vertexColors: THREE.FaceColors,
            morphTargets: false
        }));
        scene.add(points);

        /*var globeTex = new THREE.DataTexture(globeTexData, worldMap.width, worldMap.height, THREE.RGBFormat);
        globeTex.needsUpdate = true;
        globeTex.flipY = true;*/
        geometry = new THREE.SphereGeometry(200, 40, 30);
        var material = new THREE.MeshBasicMaterial({
            //wireframe: true // cool effect
            color: 0xffffff,
            //map: globeTex
        });
        globeMesh = new THREE.Mesh(geometry, material);
        globeMesh.rotation.y = Math.PI;
        scene.add(globeMesh);

        // Add the Places
        var loader = new THREE.TextureLoader();
        for (i = 0; i < opts.places.length; i++) {
            var place = opts.places[i];
            var texture = loader.load(place.coverUrl);
            texture.minFilter = THREE.LinearFilter;
            var material = new THREE.PointsMaterial({
                size: 15,
                map: texture,
                alphaTest: 0.5,
                transparent: true
            });
            var particleGeo = new THREE.Geometry();
            var outVertex = new THREE.Vector3();
            convertLatLng(place.lat, place.lng, outVertex, 205);
            particleGeo.vertices.push(outVertex);
            scene.add(new THREE.Points(particleGeo, material));
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
        point.lookAt(originPos);
        point.scale.x = xy;
        point.scale.y = xy;
        point.scale.z = h; // avoid non-invertible matrix

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
            if (tts.length == 2) {// pinch to zoom gesture
                var xx = tts[0].clientX - tts[1].clientX;
                var yy = tts[0].clientY - tts[1].clientY;
                zoomTouchSpread = Math.sqrt(xx * xx + yy * yy);
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
                var zoomDelta = Math.sqrt(xx * xx + yy * yy) - zoomTouchSpread;
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
        if (mouseOver) {
            zoom(-event.deltaY * 0.3);
        }
        return false;
    }

    function onDocumentKeyDown(event) {
        var key = event.which || event.keyCode || 0;
        switch (key) {
            /*case 38: // up arrow
                zoom(100);
                event.preventDefault();
                break;
            case 40: // down arrow
                zoom(-100);
                event.preventDefault();
                break;*/
            case 38: // up arrow
            //case 87: // w
                target.y += 0.1;
                event.preventDefault();
                break;
            case 40: // down arrow
            //case 83: // s
                target.y -= 0.1;
                event.preventDefault();
                break;
            //case 65: // a
            case 37: //left arrow
                target.x -= 0.1;
                event.preventDefault();
                break;
            //case 68: // d
            case 39: //right arrow
                target.x += 0.1;
                event.preventDefault();
                break;
        }
        target.y = target.y > PI_HALF ? PI_HALF : target.y;
        target.y = target.y < -PI_HALF ? -PI_HALF : target.y;
    }

    function onWindowResize(event) { // TODO fix window resize
        var w = container.clientWidth || 600;
        var h = Math.min(600, w) || 600;
        camera.aspect = w / h;
        camera.updateProjectionMatrix();
        renderer.setSize(w, h);
    }

    function zoom(delta) {
        distanceTarget -= delta;
        distanceTarget = distanceTarget > 1000 ? 1000 : distanceTarget;
        distanceTarget = distanceTarget < 275 ? 275 : distanceTarget;
    }

    function render() {
        zoom(0); // modifies distanceTarget

        rotation.x += (target.x - rotation.x) * 0.1;
        rotation.y += (target.y - rotation.y) * 0.1;
        distance += (distanceTarget - distance) * 0.2;

        camera.position.x = distance * Math.sin(rotation.x) * Math.cos(rotation.y);
        camera.position.y = distance * Math.sin(rotation.y);
        camera.position.z = distance * Math.cos(rotation.x) * Math.cos(rotation.y);
        camera.lookAt(originPos);

        renderer.render(scene, camera);
        requestAnimationFrame(render);
    }

    init();
    this.render = render;
    this.renderer = renderer;
    this.canvas = renderer.domElement;
    this.showLatLng = function(lat, lng, zoomTarget) { // No idea why this works
        target.x = lng * Math.PI / 180 - PI_HALF;
        target.y = lat * Math.PI / 180;
        distanceTarget = zoomTarget;
    };

    return this;
};
