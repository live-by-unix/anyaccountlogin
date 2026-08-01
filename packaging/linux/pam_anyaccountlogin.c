/*
 * PAM module for AnyAccountLogin
 * 
 * This PAM module provides authentication using AnyAccountLogin's
 * flash drive and password authentication system.
 */

#include <security/pam_appl.h>
#include <security/pam_modules.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <unistd.h>

#define ANYACCOUNTLOGIN_DAEMON "/usr/local/bin/anyaccountlogin-daemon"
#define AUTH_CONFIG_FILE "/etc/anyaccountlogin/config"

PAM_EXTERN int pam_sm_authenticate(pam_handle_t *pamh, int flags, int argc, const char **argv) {
    const char *user;
    const char *flash_drive = NULL;
    const char *password = NULL;
    int retval;

    /* Get the username */
    retval = pam_get_user(pamh, &user, NULL);
    if (retval != PAM_SUCCESS) {
        return retval;
    }

    /* Parse module options */
    for (int i = 0; i < argc; i++) {
        if (strncmp(argv[i], "flash_drive=", 12) == 0) {
            flash_drive = argv[i] + 12;
        }
    }

    /* If flash drive not specified, try to auto-detect */
    if (flash_drive == NULL) {
        /* Check common mount points */
        const char *mount_points[] = {
            "/media/usb",
            "/mnt/usb",
            "/run/media/*/USB",
            NULL
        };

        for (int i = 0; mount_points[i] != NULL; i++) {
            struct stat st;
            if (stat(mount_points[i], &st) == 0 && S_ISDIR(st.st_mode)) {
                flash_drive = mount_points[i];
                break;
            }
        }
    }

    if (flash_drive == NULL) {
        return PAM_AUTH_ERR;
    }

    /* Get password from PAM */
    retval = pam_get_authtok(pamh, PAM_AUTHTOK, (const char **)&password, NULL);
    if (retval != PAM_SUCCESS) {
        return retval;
    }

    /* Validate using AnyAccountLogin daemon */
    char cmd[1024];
    snprintf(cmd, sizeof(cmd), "%s validate --flash-drive %s --user %s",
             ANYACCOUNTLOGIN_DAEMON, flash_drive, user);

    FILE *fp = popen(cmd, "w");
    if (fp == NULL) {
        return PAM_SYSTEM_ERR;
    }

    fprintf(fp, "%s", password);
    int status = pclose(fp);

    if (WIFEXITED(status) && WEXITSTATUS(status) == 0) {
        return PAM_SUCCESS;
    }

    return PAM_AUTH_ERR;
}

PAM_EXTERN int pam_sm_setcred(pam_handle_t *pamh, int flags, int argc, const char **argv) {
    return PAM_SUCCESS;
}

PAM_EXTERN int pam_sm_acct_mgmt(pam_handle_t *pamh, int flags, int argc, const char **argv) {
    return PAM_SUCCESS;
}

PAM_EXTERN int pam_sm_open_session(pam_handle_t *pamh, int flags, int argc, const char **argv) {
    return PAM_SUCCESS;
}

PAM_EXTERN int pam_sm_close_session(pam_handle_t *pamh, int flags, int argc, const char **argv) {
    return PAM_SUCCESS;
}

PAM_EXTERN int pam_sm_chauthtok(pam_handle_t *pamh, int flags, int argc, const char **argv) {
    return PAM_SUCCESS;
}
