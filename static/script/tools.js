async function changeUserName(userID, newUserName) {
    const response = await fetch("/api/changeUserName", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({"user_id": userID,
            "new_user_name": newUserName}),
    });
    if (response.ok) {
        console.log('Username updated successfully');
    } else if (response.status === 409) {
        console.log("Username already taken!!!");
    } else {
        console.log("Name change error: " + response.status);
    }
}

async function changeUserPassword(userID, newUserPassword) {
    const response = await fetch("/api/changeUserPassword", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({"user_id": userID,
            "new_user_password": newUserPassword}),
    });
    if (response.ok) {
        console.log('Password updated successfully');
    } else {
        console.log("Password change error: " + response.status);
    }
}