<template>
    <div class="cropper">
        <div class="cropper-body">
            <div class="cropper-body-left">
                <div class="img-wrp">
                    <img :src="imgURL" alt="" ref="originImage">
                    <div class="showPosition" ref="showPosition" :style="`top: ${originHeightOffset}px;
                                left: ${originWidthOffset}px;
                                width: ${originShowWidth}px;
                                height: ${originShowHeight}px;`">
                        <!-- 裁剪框 -->
                        <div class="crop-box" ref="cropBox" :style="`width: ${boxWidth}px;
                                    height: ${boxHeight}px;
                                    top: ${boxTop}px;
                                    left: ${boxLeft}px;`">
                            <!-- 中间移动块 -->
                            <div class="crop-box-move" ref="cropBoxMove"></div>
                            <!-- 内侧左蒙版 -->
                            <div class="crop-box-marking-left"></div>
                            <!-- 内侧右蒙版 -->
                            <div class="crop-box-marking-right"></div>
                            <!-- 内侧上蒙版 不可见 -->
                            <div class="crop-box-marking-top"></div>
                            <!-- 内侧下蒙版 不可见 -->
                            <div class="crop-box-marking-bottom"></div>
                            <!-- 左上控块 -->
                            <div class="cropper-point point-lt" @mousedown="(e) => startResize('lt', e)"
                                @touchstart="(e) => startResize('lt', e.touches[0])"></div>
                            <!-- 右上控块 -->
                            <div class="cropper-point point-rt" @mousedown="(e) => startResize('rt', e)"
                                @touchstart="(e) => startResize('rt', e.touches[0])"></div>
                            <!-- 左下控块 -->
                            <div class="cropper-point point-lb" @mousedown="(e) => startResize('lb', e)"
                                @touchstart="(e) => startResize('lb', e.touches[0])"></div>
                            <!-- 右下控块 -->
                            <div class="cropper-point point-rb" @mousedown="(e) => startResize('rb', e)"
                                @touchstart="(e) => startResize('rb', e.touches[0])"></div>
                        </div>
                        <!-- 上蒙版 -->
                        <div class="marking marking-top" :style="`width: ${boxWidth}px;
                                    height: ${boxTop}px;
                                    top: 0px;
                                    left: ${boxLeft}px;`"></div>
                        <!-- 右蒙版 -->
                        <div class="marking marking-right" :style="`width: ${originShowWidth - (boxLeft + boxWidth)}px;
                                    height: ${originShowHeight}px;
                                    top: 0px;
                                    left: ${boxLeft + boxWidth}px;`"></div>
                        <!-- 下蒙版 -->
                        <div class="marking marking-bottom" :style="`width: ${boxWidth}px;
                                    height: ${originShowHeight - (boxTop + boxHeight)}px;
                                    top: ${boxTop + boxHeight}px;
                                    left: ${boxLeft}px;`"></div>
                        <!-- 左蒙版 -->
                        <div class="marking marking-left" :style="`width: ${boxLeft}px;
                                    height: ${originShowHeight}px;
                                    top: 0px;
                                    left: 0px;`"></div>
                    </div>
                </div>
                <div class="tips-wrp">
                    <div>当前分辨率 {{ cutWidth }}*{{ cutHeight }}</div>
                    <div>*部分情况下您的封面将以4:3比例展示</div>
                </div>
            </div>
            <div class="cropper-body-right">
                <div class="preview-item">
                    <div class="preview-title">16:9效果预览</div>
                    <div class="preview-card">
                        <div class="pic169">
                            <img :src="imgURL" alt=""
                                :style="`width: ${prev169Width}px;
                                        height: ${prev169Height}px;
                                        transform: translateX(-${prev169OffsetX}px) translateY(-${prev169OffsetY}px);`">
                        </div>
                        <div class="title"><span>{{ title }}</span></div>
                    </div>
                </div>
                <div class="preview-item">
                    <div class="preview-title">4:3效果预览</div>
                    <div class="preview-card">
                        <div class="pic43">
                            <img :src="imgURL" alt="" :style="`width: ${prev43Width}px;
                                        height: ${prev43Height}px;
                                        transform: translateX(-${prev43OffsetX}px) translateY(-${prev43OffsetY}px);`">
                        </div>
                        <div class="title"><span>{{ title }}</span></div>
                    </div>
                </div>
            </div>
        </div>
        <canvas ref="cropperCanvas" style="display: none"></canvas>
    </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';

const props = defineProps({
    imgURL: {
        type: String,
        default: () => null,
    },
    cutOrder: {
        type: Boolean,
        default: () => false,
    },
    title: {
        type: String,
        default: () => "title",
    }
});

const emit = defineEmits(['update:cutOrder', 'cut']);

const width = ref(0);
const height = ref(0);
const boxWidth = ref(0);
const boxHeight = ref(0);
const boxTop = ref(0);
const boxLeft = ref(0);
const isResizing = ref(false);
const resizeDirection = ref('');
const startX = ref(0);
const startBoxWidth = ref(0);
const startBoxHeight = ref(0);
const startBoxLeft = ref(0);
const startBoxTop = ref(0);

const originAspectRatio = computed(() => width.value / height.value);
const originShowWidth = computed(() => originAspectRatio.value < 16 / 9 ? 270 * originAspectRatio.value : 480);
const originShowHeight = computed(() => originAspectRatio.value > 16 / 9 ? 480 / originAspectRatio.value : 270);
const originWidthOffset = computed(() => originAspectRatio.value < 16 / 9 ? (480 - originShowWidth.value) / 2 : 0);
const originHeightOffset = computed(() => originAspectRatio.value > 16 / 9 ? (270 - originShowHeight.value) / 2 : 0);
const cutOffsetX = computed(() => Math.floor((boxLeft.value / originShowWidth.value) * width.value));
const cutOffsetY = computed(() => Math.floor((boxTop.value / originShowHeight.value) * height.value));
const cutWidth = computed(() => Math.floor((boxWidth.value / originShowWidth.value) * width.value));
const cutHeight = computed(() => Math.floor((boxHeight.value / originShowHeight.value) * height.value));
const prev169Width = computed(() => 224 * (originShowWidth.value / boxWidth.value));
const prev169Height = computed(() => 126 * (originShowHeight.value / boxHeight.value));
const prev169OffsetX = computed(() => prev169Width.value * (boxLeft.value / originShowWidth.value));
const prev169OffsetY = computed(() => prev169Height.value * (boxTop.value / originShowHeight.value));
const prev43Width = computed(() => 224 * (originShowWidth.value / (boxWidth.value * 0.75)));
const prev43Height = computed(() => 168 * (originShowHeight.value / boxHeight.value));
const prev43OffsetX = computed(() => prev43Width.value * ((boxWidth.value * 0.125 + boxLeft.value) / originShowWidth.value));
const prev43OffsetY = computed(() => prev43Height.value * (boxTop.value / originShowHeight.value));

const initBox = () => {
    if (originAspectRatio.value > 16 / 9) {
        boxWidth.value = 0.9 * originShowHeight.value * (16 / 9);
        boxHeight.value = 0.9 * originShowHeight.value;
        boxTop.value = (270 - boxHeight.value) / 2;
        boxLeft.value = (originShowWidth.value - boxWidth.value) / 2;
    } else {
        boxWidth.value = 0.9 * originShowWidth.value;
        boxHeight.value = 0.9 * originShowWidth.value / (16 / 9);
        boxTop.value = (480 - boxHeight.value) / 2;
        boxLeft.value = (originShowWidth.value - boxWidth.value) / 2;
    }
};

const cropBox = ref(null);
const cropBoxMove = ref(null);
const showPosition = ref(null);
const initDrag = () => {

    let isDragging = false;
    let offsetX, offsetY;

    cropBoxMove.value.addEventListener("mousedown", (e) => {
        isDragging = true;
        offsetX = e.clientX - cropBox.value.getBoundingClientRect().left;
        offsetY = e.clientY - cropBox.value.getBoundingClientRect().top;
    });

    cropBoxMove.value.addEventListener("touchstart", (e) => {
        isDragging = true;
        offsetX = e.touches[0].clientX - cropBox.value.getBoundingClientRect().left;
        offsetY = e.touches[0].clientY - cropBox.value.getBoundingClientRect().top;
    });

    document.addEventListener("mousemove", (e) => {
        if (!isDragging) return;

        let left = e.clientX - offsetX - showPosition.value.getBoundingClientRect().left;
        let top = e.clientY - offsetY - showPosition.value.getBoundingClientRect().top;

        left = Math.max(0, left);
        left = Math.min(showPosition.offsetWidth - cropBox.offsetWidth, left);
        top = Math.max(0, top);
        top = Math.min(showPosition.offsetHeight - cropBox.offsetHeight, top);

        boxLeft.value = left;
        boxTop.value = top;
    });


    document.addEventListener("touchmove", (e) => {
        if (!isDragging) return;

        let left = e.touches[0].clientX - offsetX - showPosition.getBoundingClientRect().left;
        let top = e.touches[0].clientY - offsetY - showPosition.getBoundingClientRect().top;

        left = Math.max(0, left);
        left = Math.min(showPosition.offsetWidth - cropBox.offsetWidth, left);
        top = Math.max(0, top);
        top = Math.min(showPosition.offsetHeight - cropBox.offsetHeight, top);

        boxLeft.value = left;
        boxTop.value = top;
    }, { passive: false });

    document.addEventListener("mouseup", () => {
        isDragging = false;
    });

    document.addEventListener("touchend", () => {
        isDragging = false;
    });
};

const updateOrigin = () => {
    const img = new Image();
    img.src = props.imgURL;
    img.onload = () => {
        width.value = img.width;
        height.value = img.height;
        setTimeout(() => {
            initBox();
        }, 1);
    }
};

const startResize = (direction, event) => {
    isResizing.value = true;
    resizeDirection.value = direction; // ✅ 直接赋值给 ref
    startX.value = event.clientX;
    startBoxWidth.value = boxWidth.value;
    startBoxHeight.value = boxHeight.value;
    startBoxLeft.value = boxLeft.value;
    startBoxTop.value = boxTop.value;

    window.addEventListener('mousemove', resize);
    window.addEventListener('touchmove', resize, { passive: false });
    window.addEventListener('mouseup', stopResize);
    window.addEventListener('touchend', stopResize);
};

const resize = (event) => {
    if (!isResizing.value) return;
    const clientX = (event.touches && event.touches.length > 0) ? event.touches[0].clientX : event.clientX;
    const deltaX = clientX - startX.value;
    if (resizeDirection.value == 'lt') {
        const newWidth = startBoxWidth.value - deltaX;
        const newLeft = startBoxLeft.value + deltaX;
        const newHeight = newWidth / (16 / 9);
        const newTop = startBoxTop.value + (startBoxHeight.value - newHeight);
        if (newWidth >= 64 && newWidth <= originShowWidth.value && newLeft >= 0 && newTop >= 0) {
            boxWidth.value = newWidth;
            boxHeight.value = newHeight;
            boxLeft.value = newLeft;
            boxTop.value = newTop;
        }
    }
    else if (resizeDirection.value == 'rt') {
        const newWidth = startBoxWidth.value + deltaX;
        const newHeight = newWidth / (16 / 9);
        const newTop = startBoxTop.value + (startBoxHeight.value - newHeight);
        if (newWidth >= 64 && newWidth + startBoxLeft.value <= originShowWidth.value && newTop >= 0) {
            boxWidth.value = newWidth;
            boxHeight.value = newHeight;
            boxTop.value = newTop;
        }
    }
    else if (resizeDirection.value == 'lb') {
        const newWidth = startBoxWidth.value - deltaX;
        const newLeft = startBoxLeft.value + deltaX;
        const newHeight = newWidth / (16 / 9);
        if (newWidth >= 64 && newWidth <= originShowWidth.value && newLeft >= 0 && startBoxTop.value + newHeight <= originShowHeight.value) {
            boxWidth.value = newWidth;
            boxHeight.value = newHeight;
            boxLeft.value = newLeft;
        }
    }
    else if (resizeDirection.value == 'rb') {
        const newWidth = startBoxWidth.value + deltaX;
        const newHeight = newWidth / (16 / 9);
        if (newWidth >= 64 && newWidth + startBoxLeft.value <= originShowWidth.value && startBoxTop.value + newHeight <= originShowHeight.value) {
            boxWidth.value = newWidth;
            boxHeight.value = newHeight;
        }
    }
    if (!(event.touches && event.touches.length > 0)) {
        event.preventDefault();
    }
};

const stopResize = () => {
    window.removeEventListener('mousemove', resize);
    window.removeEventListener('touchmove', resize);
    window.removeEventListener('mouseup', stopResize);
    window.removeEventListener('touchend', stopResize);
    isResizing.value = false;
    resizeDirection.value = '';
};

const cropperCanvas = ref(null);
const cut = () => {
    const img = new Image();
    img.src = props.imgURL;
    const canvas = cropperCanvas.value;
    canvas.width = cutWidth.value;
    canvas.height = cutHeight.value;
    const ctx = canvas.getContext("2d");
    img.onload = () => {
        ctx.drawImage(img, cutOffsetX.value, cutOffsetY.value, cutWidth.value, cutHeight.value, 0, 0, canvas.width, canvas.height);
        emit("update:cutOrder", false);
        emit("cut", canvas.toDataURL('image/jpeg'));
    };
};

onMounted(() => {
    updateOrigin();
    initDrag();
});

watch(() => props.imgURL, (current) => {
    if (current) {
        updateOrigin();
    }
});

watch(() => props.cutOrder, (current) => {
    if (current) {
        cut();
    }
});
</script>

<style scoped>
.cropper {
    user-select: none;
    -webkit-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
}

.cropper-body {
    display: flex;
    position: relative;
}

.cropper-body-left {
    padding-top: 22px;
    width: 480px;
}

.img-wrp {
    height: 270px;
    width: 100%;
    border-radius: 4px;
    background-color: #000;
    overflow: hidden;
    text-align: center;
    position: relative;
}

.img-wrp>img {
    width: 100%;
    height: 100%;
    object-fit: scale-down;
}

.showPosition {
    position: absolute;
}

.crop-box {
    position: relative;
    border: 1px solid #fff;
    box-sizing: border-box;
}

.crop-box-move {
    z-index: 9;
    position: absolute;
    cursor: move;
    width: 100%;
    height: 100%;
    left: 0;
    top: 0;
}

.crop-box-marking-left {
    z-index: 8;
    position: absolute;
    width: 12.5%;
    height: 100%;
    background: rgba(0, 0, 0, .2);
    border-right: 1px dashed rgba(255, 111, 111, .6);
    left: 0;
    top: 0;
}

.crop-box-marking-right {
    z-index: 8;
    position: absolute;
    width: 12.5%;
    height: 100%;
    background: rgba(0, 0, 0, .2);
    border-left: 1px dashed rgba(255, 111, 111, .6);
    right: 0;
    top: 0;
}

.crop-box-marking-top {
    z-index: 7;
    position: absolute;
    height: 12.5%;
    width: 100%;
    left: 0;
    top: 0;
}

.crop-box-marking-bottom {
    z-index: 7;
    position: absolute;
    height: 12.5%;
    width: 100%;
    left: 0;
    bottom: 0;
}

.cropper-point {
    position: absolute;
    background-color: #fff;
    height: 12.5px;
    opacity: 1;
    width: 12.5px;
    border-radius: 50%;
    z-index: 10;
}

.point-lt {
    cursor: nwse-resize;
    left: -6px;
    top: -6px;
}

.point-rt {
    cursor: nesw-resize;
    right: -6px;
    top: -6px;
}

.point-lb {
    cursor: nesw-resize;
    left: -6px;
    bottom: -6px;
}

.point-rb {
    cursor: nwse-resize;
    right: -6px;
    bottom: -6px;
}

.marking {
    position: absolute;
    background-color: rgba(0, 0, 0, .5);
}

.tips-wrp {
    font-weight: 400;
    font-size: 12px;
    line-height: 20px;
    color: var(--text2);
    margin-top: 12px;
}

.cropper-body-right {
    position: absolute;
    right: 0;
    transform: scale(.62);
    transform-origin: right 0;
    padding-right: 50px;
}

.preview-item {
    margin-top: 18px;
}

.preview-title {
    margin-bottom: 10px;
    position: relative;
    font-weight: 400;
    text-align: left;
    font-size: 20px;
    color: #19181c;
}

.preview-card {
    width: 224px;
    background: #fff;
    box-shadow: 0 2px 5px rgba(0, 0, 0, .1);
    border-radius: 8px;
    display: flex;
    flex-direction: column;
    align-items: center;
    position: relative;
}

.pic169 {
    position: relative;
    overflow: hidden;
    height: 126px;
    width: 100%;
    border-radius: 8px 8px 0 0;
    background-color: var(--graph_bg_regular);
}

.pic43 {
    position: relative;
    overflow: hidden;
    height: 168px;
    width: 100%;
    border-radius: 8px 8px 0 0;
    background-color: var(--graph_bg_regular);
}

.title {
    height: 65px;
    width: 100%;
}

.title span {
    margin: 8px 12px;
    font-weight: 400;
    font-size: 16px;
    line-height: 24px;
    color: #212121;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    overflow: hidden;
    -webkit-line-clamp: 2;
    word-break: break-all;
}
</style>