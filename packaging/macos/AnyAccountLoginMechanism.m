/*
 * AnyAccountLogin LoginWindow Mechanism
 * 
 * This is a loginwindow plugin for macOS that provides AnyAccountLogin
 * authentication at the login screen.
 */

#import <Cocoa/Cocoa.h>
#import <SecurityAuthorization/Authorization.h>
#import <SecurityAuthorization/AuthorizationTags.h>

@interface AnyAccountLoginMechanism : NSObject {
    AuthorizationRef authorizationRef;
}

- (BOOL)authenticateWithFlashDrive:(NSString *)flashDrivePath password:(NSString *)password;
- (void)setupLoginUI;
- (void)cleanup;

@end

@implementation AnyAccountLoginMechanism

- (id)init {
    self = [super init];
    if (self) {
        OSStatus status = AuthorizationCreate(NULL, NULL, kAuthorizationFlagDefaults, &authorizationRef);
        if (status != errAuthorizationSuccess) {
            NSLog(@"Failed to create authorization reference");
        }
    }
    return self;
}

- (void)dealloc {
    [self cleanup];
    [super dealloc];
}

- (void)cleanup {
    if (authorizationRef != NULL) {
        AuthorizationFree(authorizationRef, kAuthorizationFlagDefaults);
        authorizationRef = NULL;
    }
}

- (BOOL)authenticateWithFlashDrive:(NSString *)flashDrivePath password:(NSString *)password {
    // Validate flash drive exists
    NSFileManager *fileManager = [NSFileManager defaultManager];
    if (![fileManager fileExistsAtPath:flashDrivePath]) {
        NSLog(@"Flash drive path does not exist: %@", flashDrivePath);
        return NO;
    }

    // Check for PasswordAuthCode.txt
    NSString *authCodePath = [flashDrivePath stringByAppendingPathComponent:@"PasswordAuthCode.txt"];
    if (![fileManager fileExistsAtPath:authCodePath]) {
        NSLog(@"PasswordAuthCode.txt not found on flash drive");
        return NO;
    }

    // Read and validate auth code
    NSError *error = nil;
    NSString *storedAuthCode = [NSString stringWithContentsOfFile:authCodePath 
                                                         encoding:NSUTF8StringEncoding 
                                                            error:&error];
    if (error) {
        NSLog(@"Failed to read auth code: %@", error);
        return NO;
    }

    // Compare passwords
    if (![password isEqualToString:[storedAuthCode stringByTrimmingCharactersInSet:
                                    [NSCharacterSet whitespaceAndNewlineCharacterSet]]]) {
        NSLog(@"Password mismatch");
        return NO;
    }

    // Check for PEM key
    NSString *pemPath = [flashDrivePath stringByAppendingPathComponent:@"PasswordAuth.pem"];
    if (![fileManager fileExistsAtPath:pemPath]) {
        NSLog(@"PasswordAuth.pem not found on flash drive");
        return NO;
    }

    return YES;
}

- (void)setupLoginUI {
    // Create login window UI
    NSWindow *loginWindow = [[NSWindow alloc] initWithContentRect:NSMakeRect(0, 0, 400, 300)
                                                        styleMask:NSTitledWindowMask
                                                          backing:NSBackingStoreBuffered
                                                            defer:NO];
    
    [loginWindow setTitle:@"AnyAccountLogin"];
    [loginWindow center];
    
    // Create flash drive path field
    NSTextField *flashDriveLabel = [[NSTextField alloc] initWithFrame:NSMakeRect(20, 220, 360, 24)];
    [flashDriveLabel setStringValue:@"Flash Drive Path:"];
    [flashDriveLabel setBezeled:NO];
    [flashDriveLabel setDrawsBackground:NO];
    [flashDriveLabel setEditable:NO];
    [flashDriveLabel setSelectable:NO];
    [[loginWindow contentView] addSubview:flashDriveLabel];
    
    NSTextField *flashDriveField = [[NSTextField alloc] initWithFrame:NSMakeRect(20, 190, 360, 24)];
    [flashDriveField setStringValue:@"/Volumes/USB"];
    [[loginWindow contentView] addSubview:flashDriveField];
    
    // Create password field
    NSTextField *passwordLabel = [[NSTextField alloc] initWithFrame:NSMakeRect(20, 160, 360, 24)];
    [passwordLabel setStringValue:@"Password:"];
    [passwordLabel setBezeled:NO];
    [passwordLabel setDrawsBackground:NO];
    [passwordLabel setEditable:NO];
    [passwordLabel setSelectable:NO];
    [[loginWindow contentView] addSubview:passwordLabel];
    
    NSSecureTextField *passwordField = [[NSSecureTextField alloc] initWithFrame:NSMakeRect(20, 130, 360, 24)];
    [[loginWindow contentView] addSubview:passwordField];
    
    // Create login button
    NSButton *loginButton = [[NSButton alloc] initWithFrame:NSMakeRect(150, 80, 100, 32)];
    [loginButton setTitle:@"Login"];
    [loginButton setBezelStyle:NSRoundedBezelStyle];
    [loginButton setTarget:self];
    [loginButton setAction:@selector(performLogin:)];
    [[loginWindow contentView] addSubview:loginButton];
    
    [loginWindow makeKeyAndOrderFront:nil];
}

- (void)performLogin:(id)sender {
    // Perform authentication and notify loginwindow
    // This would integrate with the actual loginwindow mechanism
    NSLog(@"Login button clicked");
}

@end
