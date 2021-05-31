import { MainScene } from './scenes';

const mainScene = new MainScene();

window.onload = () => {
  const resizeGame = () => {
    const canvas = document.querySelector('canvas');
    const windowWidth = window.innerWidth;
    const windowHeight = window.innerHeight;
    canvas.style.width = `${window.innerWidth}px`;
    canvas.style.height = `${window.innerHeight}px`;
  };

  const config: Phaser.Types.Core.GameConfig = {
    width: 800,
    height: 600,
    backgroundColor: '#666666',
    scene: [mainScene],
  };
  const game = new Phaser.Game(config);
  window.focus();
  resizeGame();
  window.addEventListener('resize', resizeGame);
};
