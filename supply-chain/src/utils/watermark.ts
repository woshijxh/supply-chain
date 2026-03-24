/**
 * 水印工具
 * 在页面生成全局水印，显示当前登录用户名
 */

interface WatermarkOptions {
  text: string
  fontSize?: number
  color?: string
  gap?: number
  opacity?: number
  rotate?: number
}

let watermarkDiv: HTMLDivElement | null = null
let watermarkObserver: MutationObserver | null = null

/**
 * 生成水印 Canvas
 */
function createWatermarkCanvas(options: WatermarkOptions): string {
  const {
    text,
    fontSize = 16,
    color = '#999',
    gap = 100,
    opacity = 0.15,
    rotate = -22
  } = options

  const canvas = document.createElement('canvas')
  const ctx = canvas.getContext('2d')
  if (!ctx) return ''

  // 计算单个水印单元大小
  const textWidth = text.length * fontSize + 40
  const textHeight = fontSize + 20

  canvas.width = textWidth + gap
  canvas.height = textHeight + gap

  ctx.font = `${fontSize}px Arial, sans-serif`
  ctx.fillStyle = color
  ctx.globalAlpha = opacity
  ctx.textAlign = 'center'
  ctx.textBaseline = 'middle'

  // 旋转绘制
  ctx.translate(canvas.width / 2, canvas.height / 2)
  ctx.rotate((rotate * Math.PI) / 180)
  ctx.fillText(text, 0, 0)

  return canvas.toDataURL('image/png')
}

/**
 * 设置水印
 */
export function setWatermark(options: WatermarkOptions) {
  // 先清除旧水印
  removeWatermark()

  const canvasUrl = createWatermarkCanvas(options)
  if (!canvasUrl) return

  // 创建水印容器
  watermarkDiv = document.createElement('div')
  watermarkDiv.id = '__watermark__'

  Object.assign(watermarkDiv.style, {
    position: 'fixed',
    top: '0',
    left: '0',
    width: '100%',
    height: '100%',
    pointerEvents: 'none',
    zIndex: '9999',
    backgroundImage: `url(${canvasUrl})`,
    backgroundRepeat: 'repeat',
    backgroundPosition: '0 0'
  })

  document.body.appendChild(watermarkDiv)

  // 监听水印元素变化，防止被删除
  watermarkObserver = new MutationObserver((mutations) => {
    for (const mutation of mutations) {
      // 水印被删除，重新添加
      if (mutation.type === 'childList') {
        const removed = Array.from(mutation.removedNodes)
        if (removed.some(node => node === watermarkDiv)) {
          if (watermarkDiv && !document.body.contains(watermarkDiv)) {
            document.body.appendChild(watermarkDiv)
          }
        }
      }
      // 水印属性被修改，恢复
      if (mutation.type === 'attributes' && mutation.target === watermarkDiv) {
        watermarkDiv.style.display = 'block'
        watermarkDiv.style.visibility = 'visible'
        watermarkDiv.style.opacity = '1'
      }
    }
  })

  watermarkObserver.observe(document.body, {
    childList: true,
    subtree: true,
    attributes: true,
    attributeFilter: ['style', 'class']
  })
}

/**
 * 移除水印
 */
export function removeWatermark() {
  if (watermarkObserver) {
    watermarkObserver.disconnect()
    watermarkObserver = null
  }

  if (watermarkDiv && watermarkDiv.parentNode) {
    watermarkDiv.parentNode.removeChild(watermarkDiv)
    watermarkDiv = null
  }
}

/**
 * 更新水印内容
 */
export function updateWatermark(text: string) {
  if (watermarkDiv) {
    const canvasUrl = createWatermarkCanvas({ text })
    if (canvasUrl) {
      watermarkDiv.style.backgroundImage = `url(${canvasUrl})`
    }
  } else {
    setWatermark({ text })
  }
}