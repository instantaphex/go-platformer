#!/usr/bin/env node
const animations = require(process.argv[2]);
const fs = require('fs');
const assets = {};

Object.keys(animations.sprites).forEach(key => {
  const split = key.split('/');
  let newKey, animNum, isAnimation;
  if (split.length > 1) {
    animNum = split.pop();
    newKey = split.join('/');
    if (animNum.length > 2) {
      // we're not dealing with an animation
      newKey = key;
    } else {
      isAnimation = true;
    }
  } else {
    newKey = key;
  }
  if (!assets[newKey]) {
    assets[newKey] = [];
  }
  if (isAnimation) {
    assets[newKey][parseInt(animNum) - 1] = animations.sprites[key]
  } else {
    assets[newKey].push(animations.sprites[key]);
  }
});
// console.log(assets)
fs.writeFileSync('assets.json', JSON.stringify({frames:assets}, null, 2));
