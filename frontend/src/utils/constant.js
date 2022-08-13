export const options = {
  edges: {
    arrows: {
      to: {
        enabled: true,
        scaleFactor: 1,
        type: "arrow"
      },
      from: {
        enabled: false,
        scaleFactor: 1,
        type: "arrow"
      }
    },
    arrowStrikethrough: true,
    chosen: true,
    // transition line color
    color: {
      color: "#848484",
      highlight: "#848484",
      hover: "#848484",
      inherit: "from",
      opacity: 1.0
    },
    dashes: false,
    hoverWidth: 1.5,
    labelHighlightBold: true,
    physics: false,
  },
  layout: {
    randomSeed: 0.1,
  }
}