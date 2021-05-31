import * as Phaser from 'phaser';
import { Session } from './session';
import squares_png from './assets/squares.png';
import font24_png from './assets/font24.png';
import font24_fnt from './assets/font24.fnt';
import font96_png from './assets/font96.png';
import font96_fnt from './assets/font96.fnt';

export class MainScene extends Phaser.Scene {
  readonly cell = 20;
  conn: Session;
  tiles: Phaser.GameObjects.Sprite[];
  rows: number;
  cols: number;
  overflow: number;
  cursors: Phaser.Types.Input.Keyboard.CursorKeys;
  scoreTxt: Phaser.GameObjects.BitmapText;
  levelTxt: Phaser.GameObjects.BitmapText;
  gameOverText: Phaser.GameObjects.BitmapText;

  upPressed: boolean;
  downPressed: boolean;
  leftPressed: boolean;
  rightPressed: boolean;
  spacePressed: boolean;

  constructor() {
    super('');
    this.conn = new Session();
  }

  preload() {
    this.load.json('meta', this.conn.metaUrl());
    this.load.spritesheet('squares', squares_png, {
      frameWidth: this.cell,
      frameHeight: this.cell,
    });
    this.load.bitmapFont('font24', font24_png, font24_fnt);
    this.load.bitmapFont('font96', font96_png, font96_fnt);
  }

  stateCange(msg: string) {
    const state = JSON.parse(msg);
    const t = state['tiles'];
    for (let r = this.overflow; r < this.rows; r++) {
      for (let c = 0; c < this.cols; c++) {
        const ndx = r * this.cols + c;
        this.tiles[ndx].setFrame(t[ndx]);
      }
    }
    const level = state['level'];
    const score = state['score'];
    const gameOver = state['gameOver'];
    this.levelTxt.text = `level: ${level}`;
    this.scoreTxt.text = `score: ${score}`;
    this.gameOverText.visible = gameOver;
  }

  create() {
    this.cursors = this.input.keyboard.createCursorKeys();
    this.conn.onmessage = (msg) => this.stateCange(msg);
    this.conn.connect();
    const meta = this.cache.json.get('meta');
    this.rows = meta['rows'];
    this.cols = meta['cols'];
    this.overflow = meta['overflow'];
    this.tiles = new Array<Phaser.GameObjects.Sprite>(this.rows * this.cols);
    for (let r = this.overflow; r < this.rows; r++) {
      for (let c = 0; c < this.cols; c++) {
        const tile = this.add
          .sprite(c * this.cell, (r - this.overflow) * this.cell, 'squares', 0)
          .setOrigin(0, 0);
        this.tiles[r * this.cols + c] = tile;
      }
    }
    this.levelTxt = this.add.bitmapText(216, 32, 'font24', '');
    this.scoreTxt = this.add.bitmapText(216, 56, 'font24', '');
    this.gameOverText = this.add.bitmapText(216, 200, 'font96', 'Game Over');
    this.gameOverText.visible = false;
  }

  update() {
    //
    if (this.cursors.left.isDown) {
      if (!this.leftPressed) {
        this.conn.command('left');
      }
      this.leftPressed = true;
    } else {
      this.leftPressed = false;
    }
    //
    if (this.cursors.right.isDown) {
      if (!this.rightPressed) {
        this.conn.command('right');
      }
      this.rightPressed = true;
    } else {
      this.rightPressed = false;
    }
    //
    if (this.cursors.up.isDown) {
      if (!this.upPressed) {
        this.conn.command('up');
      }
      this.upPressed = true;
    } else {
      this.upPressed = false;
    }
    //
    if (this.cursors.down.isDown) {
      if (!this.downPressed) {
        this.conn.command('down');
      }
      this.downPressed = true;
    } else {
      this.downPressed = false;
    }
    //
    if (this.cursors.space.isDown) {
      if (!this.spacePressed) {
        this.conn.command('space');
      }
      this.spacePressed = true;
    } else {
      this.spacePressed = false;
    }
  }
}
