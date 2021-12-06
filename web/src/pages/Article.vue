<template>
    <h2 class="text-center">online edit</h2>
    <div :id="tinymceId"></div>
    <div>
        <button @click="save" :class="classEditMode">save</button>
    </div>
</template>

<script>
    import {toRef, unref} from 'vue'

    import jQuery from 'jquery'

    import Editor from '@tinymce/tinymce-vue'

    /* Import TinyMCE */
    import tinymce from 'tinymce';

    /* Default icons are required for TinyMCE 5.3 or above */
    import 'tinymce/icons/default';

    /* A theme is also required */
    import 'tinymce/themes/silver';

    /* Import the skin */
    import 'tinymce/skins/ui/oxide/skin.css';

    /* Import plugins */
    import 'tinymce/plugins/advlist';
    import 'tinymce/plugins/code';
    import 'tinymce/plugins/emoticons';
    import 'tinymce/plugins/emoticons/js/emojis';
    import 'tinymce/plugins/charmap';
    import 'tinymce/plugins/lists';
    import 'tinymce/plugins/link';
    import 'tinymce/plugins/table';
    import 'tinymce/plugins/wordcount';
    import 'tinymce/plugins/image';
    import 'tinymce/plugins/imagetools';
    import 'tinymce/plugins/anchor';
    import 'tinymce/plugins/preview';
    import 'tinymce/plugins/visualblocks';
    import 'tinymce/plugins/fullscreen';
    import 'tinymce/plugins/insertdatetime';
    import 'tinymce/plugins/media';
    import 'tinymce/plugins/paste';

    /* Import content css */
    import 'tinymce/skins/ui/oxide/content.css';
    import 'tinymce/skins/content/default/content.css';

    export default {
        name: "Article",
        props: {
            articleKey: String,
            suppressDataUpdate: {
                type: Boolean,
                default: false,
            },
            suppressEdit: {
                type: Boolean,
                default: true,
            },
        },
        setup(props) {
            const articleKey = toRef(props, 'articleKey')
            const suppressDataUpdate = toRef(props, 'suppressDataUpdate')
            const suppressEdit = toRef(props, 'suppressEdit')

            return {
                meArticleKey: unref(articleKey),
                meSuppressDataUpdate: unref(suppressDataUpdate),
                meSuppressEdit: unref(suppressEdit),
            }
        },
        data() {
            const id = "tiny-id-" + Date.now()
            return {
                tinymceId: id,
                tinymceEditor: '',
                defaultInit: {
                    height: 500,
                    skin: false,
                    content_css: '/static/tinymce_content.css',
                    image_dimensions: false,
                    image_title: false,
                    image_description: false,
                    paste_data_images: true,
                    object_resizing: 'img',
                    imagetools_toolbar: 'rotateleft rotateright | flipv fliph | editimage',
                    contextmenu: 'link',
                    plugins: [
                        'advlist lists link image imagetools charmap preview anchor',
                        'visualblocks code fullscreen',
                        'insertdatetime media table paste wordcount'
                    ],
                    toolbar:
                        'undo redo | bold italic backcolor | \
                        alignleft aligncenter alignright alignjustify | \
                        bullist numlist outdent indent | removeformat',
                    images_upload_handler: function (blobInfo, success, failure, progress) {
                        let formData = new FormData()
                        formData.append('target', 'smSite')
                        formData.append('file', blobInfo.blob(),)

                        jQuery.ajax({
                            url: '/api/upload',
                            data: formData,
                            dataType: 'json',
                            type: 'POST',
                            contentType: false,
                            processData: false,
                            success: function (data) {
                                console.log(data)

                                if (data && data.success) {
                                    success(data.data.url)
                                } else {
                                    failure()
                                }
                            },
                            error: function (data) {
                                console.log(data)

                                failure()
                            }
                        })
                    },
                }
            }
        },
        computed: {
            classEditMode() {
                return {
                    'not-active': this.meSuppressEdit
                }
            }
        },
        created() {
            this.$watch(
                () => this.$route,
                (toRoute, previousRoute) => {
                    if (previousRoute === undefined) return

                    this.meArticleKey = toRoute.params.key
                    this.meSuppressDataUpdate = false
                    this.meSuppressEdit = true
                    if (toRoute.params.hasOwnProperty("suppressDataUpdate")) {
                        this.meSuppressDataUpdate = toRoute.params.suppressDataUpdate === "true"
                    }
                    if (toRoute.query && toRoute.query.hasOwnProperty("mode")) {
                        this.meSuppressEdit = !(toRoute.query.mode === "edit")
                    }

                    this.fetchData()
                },
                {immediate: true}
            )
        },
        mounted() {
            this.init()
        },
        methods: {
            init() {
                const self = this

                self.initEditor()

                self.fetchData()
            },
            initEditor() {
                const self = this

                if (self.tinymceEditor) {
                    self.tinymceEditor.remove()
                }

                let customizeInit = {}
                if (self.meSuppressEdit) {
                    Object.assign(customizeInit, self.defaultInit, {
                        inline: true,
                        toolbar: false,
                        menubar: false,
                        readonly: true,
                    })
                } else {
                    Object.assign(customizeInit, self.defaultInit, {
                        menubar: false,
                    })
                }

                window.tinymce.init({
                    ...customizeInit,
                    setup: function (editor) {
                        self.tinymceEditor = editor
                    },
                    selector: `#${self.tinymceId}`
                })
            },
            save(event) {
                const self = this
                console.log(this.tinymceEditor.getContent())

                if (self.meArticleKey) {
                    jQuery.ajax({
                        url: '/update-article',
                        data: JSON.stringify({
                            key: self.meArticleKey,
                            content: self.tinymceEditor.getContent(),
                        }),
                        contentType: 'application/json',
                        type: 'POST',
                        success: function (data) {
                            console.log(data)

                            if (data && data.success) {

                            }
                        },
                        error: function (data) {
                            console.log(data)
                        }
                    })
                } else {
                    jQuery.ajax({
                        url: '/create-article',
                        data: JSON.stringify({
                            content: self.tinymceEditor.getContent(),
                        }),
                        contentType: 'application/json',
                        type: 'POST',
                        success: function (data) {
                            console.log(data)

                            if (data && data.success) {
                                self.$router.push({
                                    name: 'Article',
                                    params: {
                                        key: data.data.key,
                                        suppressDataUpdate: true,
                                    },
                                    query: {
                                        mode: "edit"
                                    }
                                })
                            }
                        },
                        error: function (data) {
                            console.log(data)
                        }
                    })
                }
            },
            fetchData() {
                let self = this

                if (!self.meArticleKey) {
                    if (!self.tinymceEditor) return

                    self.tinymceEditor.resetContent("")
                    return
                }

                if (self.meSuppressDataUpdate) return

                jQuery.ajax({
                    url: '/fetch-article',
                    data: JSON.stringify({
                        key: self.meArticleKey,
                    }),
                    contentType: 'application/json',
                    type: 'POST',
                    success: function (data) {
                        console.log(data)

                        if (data && data.success) {
                            self.tinymceEditor.resetContent(data.data.content)
                        } else {
                            self.$router.push({name: 'Home'})
                        }
                    },
                    error: function (data) {
                        console.log(data)
                    }
                })
            },
        },
    }
</script>

<style scoped>
    .text-center {
        text-align: center;
    }

    .not-active {
        display: none;
    }
</style>