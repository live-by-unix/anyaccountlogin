/*
 * AnyAccountLogin Windows Service
 * 
 * This is a Windows service that runs the AnyAccountLogin daemon
 * at system startup for authentication integration.
 */

#include <windows.h>
#include <winsvc.h>
#include <stdio.h>
#include <tchar.h>

#define SERVICE_NAME _T("AnyAccountLoginService")
#define SERVICE_DISPLAY_NAME _T("AnyAccountLogin Authentication Service")
#define SERVICE_DESCRIPTION _T("Provides AnyAccountLogin authentication services")

SERVICE_STATUS g_serviceStatus;
SERVICE_STATUS_HANDLE g_serviceStatusHandle;
HANDLE g_serviceStopEvent = NULL;

void WINAPI ServiceMain(DWORD argc, LPTSTR *argv);
void WINAPI ServiceCtrlHandler(DWORD dwCtrl);
DWORD ServiceWorkerThread(LPVOID lpParam);
BOOL InstallService();
BOOL UninstallService();
void LogEvent(LPCTSTR pszMessage, WORD wType);

int _tmain(int argc, TCHAR *argv[])
{
    if (argc > 1)
    {
        if (_tcscmp(argv[1], _T("install")) == 0)
        {
            if (InstallService())
            {
                _tprintf(_T("Service installed successfully.\n"));
                return 0;
            }
            else
            {
                _tprintf(_T("Failed to install service.\n"));
                return 1;
            }
        }
        else if (_tcscmp(argv[1], _T("uninstall")) == 0)
        {
            if (UninstallService())
            {
                _tprintf(_T("Service uninstalled successfully.\n"));
                return 0;
            }
            else
            {
                _tprintf(_T("Failed to uninstall service.\n"));
                return 1;
            }
        }
        else
        {
            _tprintf(_T("Usage: %s [install|uninstall]\n"), argv[0]);
            return 1;
        }
    }

    SERVICE_TABLE_ENTRY serviceTable[] = {
        { (LPWSTR)SERVICE_NAME, (LPSERVICE_MAIN_FUNCTION)ServiceMain },
        { NULL, NULL }
    };

    if (!StartServiceCtrlDispatcher(serviceTable))
    {
        LogEvent(_T("Failed to start service control dispatcher"), EVENTLOG_ERROR_TYPE);
        return 1;
    }

    return 0;
}

void WINAPI ServiceMain(DWORD argc, LPTSTR *argv)
{
    g_serviceStatusHandle = RegisterServiceCtrlHandler(SERVICE_NAME, ServiceCtrlHandler);
    if (g_serviceStatusHandle == NULL)
    {
        LogEvent(_T("Failed to register service control handler"), EVENTLOG_ERROR_TYPE);
        return;
    }

    ZeroMemory(&g_serviceStatus, sizeof(g_serviceStatus));
    g_serviceStatus.dwServiceType = SERVICE_WIN32_OWN_PROCESS;
    g_serviceStatus.dwControlsAccepted = SERVICE_ACCEPT_STOP | SERVICE_ACCEPT_SHUTDOWN;
    g_serviceStatus.dwCurrentState = SERVICE_START_PENDING;
    g_serviceStatus.dwWin32ExitCode = 0;
    g_serviceStatus.dwServiceSpecificExitCode = 0;
    g_serviceStatus.dwCheckPoint = 0;
    g_serviceStatus.dwWaitHint = 0;

    SetServiceStatus(g_serviceStatusHandle, &g_serviceStatus);

    g_serviceStopEvent = CreateEvent(NULL, TRUE, FALSE, NULL);
    if (g_serviceStopEvent == NULL)
    {
        LogEvent(_T("Failed to create service stop event"), EVENTLOG_ERROR_TYPE);
        g_serviceStatus.dwCurrentState = SERVICE_STOPPED;
        g_serviceStatus.dwWin32ExitCode = GetLastError();
        SetServiceStatus(g_serviceStatusHandle, &g_serviceStatus);
        return;
    }

    g_serviceStatus.dwCurrentState = SERVICE_RUNNING;
    SetServiceStatus(g_serviceStatusHandle, &g_serviceStatus);

    LogEvent(_T("AnyAccountLogin Service started"), EVENTLOG_INFORMATION_TYPE);

    // Start the worker thread
    HANDLE hThread = CreateThread(NULL, 0, (LPTHREAD_START_ROUTINE)ServiceWorkerThread, NULL, 0, NULL);
    if (hThread == NULL)
    {
        LogEvent(_T("Failed to create worker thread"), EVENTLOG_ERROR_TYPE);
        CloseHandle(g_serviceStopEvent);
        g_serviceStatus.dwCurrentState = SERVICE_STOPPED;
        g_serviceStatus.dwWin32ExitCode = GetLastError();
        SetServiceStatus(g_serviceStatusHandle, &g_serviceStatus);
        return;
    }

    WaitForSingleObject(g_serviceStopEvent, INFINITE);

    CloseHandle(hThread);
    CloseHandle(g_serviceStopEvent);

    g_serviceStatus.dwCurrentState = SERVICE_STOPPED;
    SetServiceStatus(g_serviceStatusHandle, &g_serviceStatus);

    LogEvent(_T("AnyAccountLogin Service stopped"), EVENTLOG_INFORMATION_TYPE);
}

void WINAPI ServiceCtrlHandler(DWORD dwCtrl)
{
    switch (dwCtrl)
    {
    case SERVICE_CONTROL_STOP:
    case SERVICE_CONTROL_SHUTDOWN:
        g_serviceStatus.dwCurrentState = SERVICE_STOP_PENDING;
        SetServiceStatus(g_serviceStatusHandle, &g_serviceStatus);
        SetEvent(g_serviceStopEvent);
        break;
    default:
        break;
    }
}

DWORD ServiceWorkerThread(LPVOID lpParam)
{
    // This is where the Go daemon would be executed
    // For now, we'll just keep the service running
    
    while (WaitForSingleObject(g_serviceStopEvent, 1000) != WAIT_OBJECT_0)
    {
        // Service is running
        // In a real implementation, this would launch and monitor the Go daemon
    }

    return 0;
}

BOOL InstallService()
{
    SC_HANDLE hSCManager = OpenSCManager(NULL, NULL, SC_MANAGER_CREATE_SERVICE);
    if (hSCManager == NULL)
    {
        LogEvent(_T("Failed to open service control manager"), EVENTLOG_ERROR_TYPE);
        return FALSE;
    }

    TCHAR szPath[MAX_PATH];
    if (!GetModuleFileName(NULL, szPath, MAX_PATH))
    {
        LogEvent(_T("Failed to get module file name"), EVENTLOG_ERROR_TYPE);
        CloseServiceHandle(hSCManager);
        return FALSE;
    }

    SC_HANDLE hService = CreateService(
        hSCManager,
        SERVICE_NAME,
        SERVICE_DISPLAY_NAME,
        SERVICE_ALL_ACCESS,
        SERVICE_WIN32_OWN_PROCESS,
        SERVICE_AUTO_START,
        SERVICE_ERROR_NORMAL,
        szPath,
        NULL, NULL, NULL, NULL, NULL);

    if (hService == NULL)
    {
        LogEvent(_T("Failed to create service"), EVENTLOG_ERROR_TYPE);
        CloseServiceHandle(hSCManager);
        return FALSE;
    }

    // Set service description
    SERVICE_DESCRIPTION sd;
    sd.lpDescription = (LPTSTR)SERVICE_DESCRIPTION;
    ChangeServiceConfig2(hService, SERVICE_CONFIG_DESCRIPTION, &sd);

    CloseServiceHandle(hService);
    CloseServiceHandle(hSCManager);

    LogEvent(_T("Service installed successfully"), EVENTLOG_INFORMATION_TYPE);
    return TRUE;
}

BOOL UninstallService()
{
    SC_HANDLE hSCManager = OpenSCManager(NULL, NULL, SC_MANAGER_CONNECT);
    if (hSCManager == NULL)
    {
        LogEvent(_T("Failed to open service control manager"), EVENTLOG_ERROR_TYPE);
        return FALSE;
    }

    SC_HANDLE hService = OpenService(hSCManager, SERVICE_NAME, SERVICE_STOP | DELETE);
    if (hService == NULL)
    {
        LogEvent(_T("Failed to open service"), EVENTLOG_ERROR_TYPE);
        CloseServiceHandle(hSCManager);
        return FALSE;
    }

    SERVICE_STATUS status;
    ControlService(hService, SERVICE_CONTROL_STOP, &status);

    if (!DeleteService(hService))
    {
        LogEvent(_T("Failed to delete service"), EVENTLOG_ERROR_TYPE);
        CloseServiceHandle(hService);
        CloseServiceHandle(hSCManager);
        return FALSE;
    }

    CloseServiceHandle(hService);
    CloseServiceHandle(hSCManager);

    LogEvent(_T("Service uninstalled successfully"), EVENTLOG_INFORMATION_TYPE);
    return TRUE;
}

void LogEvent(LPCTSTR pszMessage, WORD wType)
{
    HANDLE hEventSource = RegisterEventSource(NULL, SERVICE_NAME);
    if (hEventSource != NULL)
    {
        LPCTSTR lpszStrings[1] = { pszMessage };
        ReportEvent(hEventSource, wType, 0, 0, NULL, 1, 0, lpszStrings, NULL);
        DeregisterEventSource(hEventSource);
    }
}
