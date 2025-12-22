//go:build android

#include <android/log.h>
#include <jni.h>
#include <stdbool.h>
#include <stdlib.h>
#include <string.h>

#define LOG_FATAL(...) __android_log_print(ANDROID_LOG_FATAL, "Fyne", __VA_ARGS__)

// helper functions go here
static jclass find_class(JNIEnv *env, const char *class_name) {
    jclass clazz = (*env)->FindClass(env, class_name);
    if (clazz == NULL) {
        (*env)->ExceptionClear(env);
        LOG_FATAL("cannot find %s", class_name);
        return NULL;
    }
    return clazz;
}

static jmethodID find_method(JNIEnv *env, jclass clazz, const char *name, const char *sig) {
    jmethodID m = (*env)->GetMethodID(env, clazz, name, sig);
    if (m == 0) {
        (*env)->ExceptionClear(env);
        LOG_FATAL("cannot find method %s %s", name, sig);
        return 0;
    }
    return m;
}

static jmethodID find_static_method(JNIEnv *env, jclass clazz, const char *name, const char *sig) {
    jmethodID m = (*env)->GetStaticMethodID(env, clazz, name, sig);
    if (m == 0) {
        (*env)->ExceptionClear(env);
        LOG_FATAL("cannot find method %s %s", name, sig);
        return 0;
    }
    return m;
}

const char* getCString(uintptr_t jni_env, uintptr_t ctx, jstring str) {
    JNIEnv *env = (JNIEnv*)jni_env;

    const char *chars = (*env)->GetStringUTFChars(env, str, NULL);

    const char *copy = strdup(chars);
    (*env)->ReleaseStringUTFChars(env, str, chars);
    return copy;
}

const char *androidName(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx) {
    JNIEnv *env = (JNIEnv*)jni_env;
    // look up model from Build class
    jclass buildClass = find_class(env, "android/os/Build");
    jfieldID modelFieldID = (*env)->GetStaticFieldID(env, buildClass, "MODEL", "Ljava/lang/String;");
    jstring model = (*env)->GetStaticObjectField(env, buildClass, modelFieldID);
    // convert to a C string
    return getCString(jni_env, ctx, model);
}
